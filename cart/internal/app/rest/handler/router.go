package handler

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"route256/cart/internal/middleware"
	"route256/cart/internal/usecase/cart"
)

type Router struct {
	cartService *cart.Service
	mux         *http.ServeMux
}

func NewRouter(cartService *cart.Service) *Router {
	return &Router{
		cartService: cartService,
	}
}

func (s *Router) GetMux() *http.ServeMux {
	mux := http.NewServeMux()

	setHandler(mux, "GET", "/user/{id}/cart/list", s.GetCart)
	setHandler(mux, "POST", "/user/{id}/cart/{sku}", s.AddCartItem)
	setHandler(mux, "POST", "/user/{id}/cart/checkout", s.CheckoutCart)
	setHandler(mux, "DELETE", "/user/{id}/cart", s.RemoveCart)
	setHandler(mux, "DELETE", "/user/{id}/cart/{sku}", s.RemoveCartItem)

	return mux
}

func setHandler(mux *http.ServeMux, method string, url string, handler middleware.ErrorWrapper) {
	otelMiddleware := otelhttp.NewMiddleware("http.handler: " + method + " " + url)

	mux.Handle(method+" "+url,
		otelMiddleware(
			middlewareSet(handler, url),
		),
	)
}

func middlewareSet(handler http.Handler, url string) http.Handler {
	return otelhttp.WithRouteTag(url,
		middleware.Logging(
			middleware.RequestCounter(
				middleware.JSON(handler), url,
			),
		),
	)
}
