package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io"
	"net/http"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
	"route256/logger/pkg/logger"
	"route256/metrics/pkg/tracer"
	"strconv"
	"time"
)

const HTTPRequestMaxAttempts = 3

type errorResponse struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
}

type request struct {
	Token string `json:"token"`
}

type tokenRequest interface {
	setToken(token string)
}

func (r *request) setToken(token string) {
	r.Token = token
}

type productRequest struct {
	request
	Sku int64 `json:"sku"`
}

type productResponse struct {
	Name  string      `json:"name" validate:"required"`
	Price model.Price `json:"price" validate:"gt=0"`
}

type portRequestError struct {
	StatusCode int
	Message    string
}

type HTTPProvider struct {
	c      http.Client
	config Config
}

func NewHTTPProvider(config Config) *HTTPProvider {
	return &HTTPProvider{
		config: config,
		c: http.Client{
			Transport: NewRetryTransport(HTTPRequestMaxAttempts),
		},
	}
}

func toProduct(response *productResponse) model.Product {
	return model.Product{
		Price: response.Price,
		Name:  response.Name,
	}
}

func validateProductResponse(response *productResponse) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validErr := validate.Struct(response)
	if validErr != nil {
		return validErr
	}

	return nil
}

var requestCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "app",
		Name:      "product_request_total_counter",
		Help:      "Total amount of product requests",
	},
	[]string{},
)

var requestHistogramVec = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "app",
		Name:      "product_request_duration_histogram",
		Buckets:   prometheus.DefBuckets,
		Help:      "Duration of product requests",
	},
	[]string{"url", "statusCode"})

func (p *HTTPProvider) Get(ctx context.Context, sku model.Sku) (model.Product, error) {

	ctx, span := tracer.BeginSpan(ctx, "adapter.product.HTTPProvider/Get")
	defer span.End()

	req := productRequest{Sku: int64(sku)}

	res, err := p.postRequest(ctx, &req, "/get_product", &productResponse{})

	if err != nil {
		e := &portRequestError{}
		if errors.As(err, e) {
			logger.Warnw(ctx, "product service", "error", err.Error())

			switch e.StatusCode {
			case http.StatusNotFound:
				return cart.EmptyProduct, cart.ErrProductSkuNotFound
			case http.StatusTooManyRequests:
				return cart.EmptyProduct, cart.ErrProductServiceIsBusy
			case 420:
				return cart.EmptyProduct, cart.ErrProductServiceIsBusy
			default:
				return cart.EmptyProduct, err
			}
		} else {
			return cart.EmptyProduct, err
		}
	}

	productRes := res.(*productResponse)

	validErr := validateProductResponse(productRes)
	if validErr != nil {
		return cart.EmptyProduct, fmt.Errorf("%w: "+validErr.Error(), cart.ErrProductServiceProblem)
	}

	return toProduct(productRes), nil
}

func (e portRequestError) Error() string {
	return e.Message
}

func (r *HTTPProvider) postRequest(_ context.Context, req tokenRequest, methodName string, dst any) (any, error) {
	req.setToken(r.config.Token)

	buf := bytes.NewBuffer(make([]byte, 0))
	enc := json.NewEncoder(buf)

	err := enc.Encode(req)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	var respErr error
	func() {
		defer func(createdAt time.Time) {
			if resp != nil {
				requestCounterVec.WithLabelValues().Inc()
				requestHistogramVec.WithLabelValues(methodName, strconv.Itoa(resp.StatusCode)).Observe(time.Since(createdAt).Seconds())
			}
		}(time.Now())
		resp, respErr = r.c.Post(r.config.Url+methodName, "application/json", buf)
	}()

	if respErr != nil {
		return nil, err
	}

	rawResp, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := &errorResponse{}
		err = json.Unmarshal(rawResp, errResp)
		if err != nil {
			return nil, err
		}

		return nil, portRequestError{StatusCode: resp.StatusCode, Message: errResp.Message}
	}

	err = json.Unmarshal(rawResp, dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
