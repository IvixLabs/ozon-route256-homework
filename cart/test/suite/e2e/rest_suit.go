package e2e

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
)

type RestSuite struct {
	suite.Suite
	ctx       context.Context
	container testcontainers.Container
	host      string
	port      string
	client    *http.Client
}

func (s *RestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.client = &http.Client{}

	dockerContextPath := os.Getenv("DOCKER_CONTEXT_PATH")
	if dockerContextPath == "" {
		panic("DOCKER_CONTEXT_PATH is required")
	}

	const AppRestPort = "8082"
	const AppRestAddr = ":" + AppRestPort

	reqContainers := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    dockerContextPath,
			Dockerfile: "build/Dockerfile",
			KeepImage:  false,
			BuildOptionsModifier: func(buildOptions *types.ImageBuildOptions) {
				buildOptions.Target = "run"
			},
		},
		Env: map[string]string{
			"APP_REST_ADDR":                 AppRestAddr,
			"APP_MODE":                      "test",
			"APP_GRACEFUL_SHUTDOWN_TIMEOUT": "5",
			"APP_DEBUG_SRV":                 "disabled",
			"APP_TRACER":                    "disabled",
		},
		ExposedPorts: []string{AppRestPort},
		WaitingFor:   wait.ForListeningPort(AppRestPort),
	}

	var err error
	s.container, err = testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: reqContainers,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	cHost, err := s.container.Host(s.ctx)
	if err != nil {
		panic(err)
	}
	s.host = cHost

	cPort, err := s.container.MappedPort(s.ctx, "8082")
	if err != nil {
		panic(err)
	}
	s.port = cPort.Port()

}

func (s *RestSuite) TearDownSuite() {
	err := s.container.Terminate(s.ctx)
	if err != nil {
		panic(err)
	}
}

func (s *RestSuite) SetupTest() {
}

func (s *RestSuite) TearDownTest() {
}
func (s *RestSuite) TestRest() {

	tests := []struct {
		name     string
		path     string
		body     string
		method   string
		wantCode int
		wantBody string
	}{
		{
			name:     "Remove_unknown_cart",
			path:     "/user/123/cart",
			method:   http.MethodDelete,
			wantCode: http.StatusNotFound,
			wantBody: "{\"message\":\"entity is not found: cart is not found\"}",
		},
		{
			name:     "Add_1_sku_to_cart",
			path:     "/user/31337/cart/111",
			method:   http.MethodPost,
			wantCode: http.StatusOK,
			body:     "{\"count\":1}",
		},
		{
			name:     "Remove_unknown_sku_from_cart",
			path:     "/user/31337/cart/123",
			method:   http.MethodDelete,
			wantCode: http.StatusNotFound,
			wantBody: "{\"message\":\"entity is not found: cart item is not found\"}",
		},
		{
			name:     "Increase_1_sku_in_cart",
			path:     "/user/31338/cart/111",
			method:   http.MethodPost,
			body:     "{\"count\":1}",
			wantCode: http.StatusOK,
		},
		{
			name:     "Increase_to_5_sku_in_cart",
			path:     "/user/31338/cart/111",
			method:   http.MethodPost,
			body:     "{\"count\":5}",
			wantCode: http.StatusOK,
		},
		{
			name:     "Add_1_unknown_sku_to_cart",
			path:     "/user/31337/cart/123",
			method:   http.MethodPost,
			wantCode: http.StatusPreconditionFailed,
			body:     "{\"count\":1}",
			wantBody: "{\"message\":\"wrong argument: product sku is not found\"}",
		},
		{
			name:     "Add_another_sku_to_cart",
			path:     "/user/31337/cart/222",
			method:   http.MethodPost,
			wantCode: http.StatusOK,
			body:     "{\"count\":1}",
		},
		{
			name:     "Add_sku_to_invalid_user",
			path:     "/user/0/cart/222",
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
			body:     "{\"count\":1}",
			wantBody: "{\"message\":\"validation error: wrong id\"}",
		},
		{
			name:     "Add_invalid_sku_to_cart",
			path:     "/user/31337/cart/0",
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
			body:     "{\"count\":1}",
			wantBody: "{\"message\":\"validation error: wrong sku\"}",
		},
		{
			name:     "Add_invalid_sku_to_cart_invalid_count",
			path:     "/user/31337/cart/111",
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
			body:     "{\"count\":0}",
			wantBody: "{\"message\":\"validation error Key: 'addCartItemRequest.Count' Error:Field validation for 'Count' failed on the 'gt' tag\"}",
		},
		{
			name:     "Remove_whole_sku_from_cart",
			path:     "/user/31337/cart/111",
			method:   http.MethodDelete,
			wantCode: http.StatusNoContent,
		},
		{
			name:     "Delete_whole_cart",
			path:     "/user/31337/cart",
			method:   http.MethodDelete,
			wantCode: http.StatusNoContent,
		},
		{
			name:     "Get_list_of_a_cart",
			path:     "/user/31338/cart/list",
			method:   http.MethodGet,
			wantCode: http.StatusOK,
			wantBody: "{\"items\":[{\"sku_id\":111,\"name\":\"Product 111\",\"count\":6,\"price\":100}],\"total_price\":600}",
		},
		{
			name:     "Get_invalid_user_cart_list",
			path:     "/user/0/cart/list",
			method:   http.MethodGet,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"message\":\"validation error: wrong id\"}",
		},
		{
			name:     "Delete_whole_cart",
			path:     "/user/31338/cart",
			method:   http.MethodDelete,
			wantCode: http.StatusNoContent,
		},
		{
			name:     "Add_1_sku_to_cart",
			path:     "/user/31337/cart/111",
			method:   http.MethodPost,
			wantCode: http.StatusOK,
			body:     "{\"count\":1}",
		},
		{
			name:     "Checkout_wrong_cart",
			path:     "/user/7777/cart/checkout",
			method:   http.MethodPost,
			wantCode: http.StatusNotFound,
			wantBody: "{\"message\":\"entity is not found: cart is not found\"}",
		},
		{
			name:     "Checkout_cart",
			path:     "/user/31337/cart/checkout",
			method:   http.MethodPost,
			wantCode: http.StatusOK,
			wantBody: "{\"orderID\":1}",
		},
		{
			name:     "Get_list_of_wrong_cart",
			path:     "/user/31337/cart/list",
			method:   http.MethodGet,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			req := s.buildRequest(tt.path, tt.method, tt.body)
			code, body := s.doRequest(req)

			assert.Equal(s.T(), code, tt.wantCode)
			assert.Equal(s.T(), body, tt.wantBody)
		})
	}
}

func (s *RestSuite) doRequest(req *http.Request) (int, string) {
	res, err := s.client.Do(req)
	assert.NoError(s.T(), err)

	return res.StatusCode, s.readBody(res.Body)
}

func (s *RestSuite) buildRequest(path string, method string, body string) *http.Request {
	req := &http.Request{
		URL:    &url.URL{Scheme: "http", Host: s.host + ":" + s.port, Path: path},
		Method: method,
	}

	var bodyReader io.ReadCloser
	if body != "" {
		bodyReader = io.NopCloser(bytes.NewReader([]byte(body)))
	}

	req.Body = bodyReader
	return req
}

func (s *RestSuite) readBody(reader io.ReadCloser) string {
	body, err := io.ReadAll(reader)
	defer reader.Close()

	assert.NoError(s.T(), err)

	return string(body)
}
