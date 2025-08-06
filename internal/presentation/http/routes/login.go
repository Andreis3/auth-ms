package routes

import (
	"net/http"

	"github.com/andreis3/customers-ms/internal/presentation/http/handler"
	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
	"github.com/andreis3/customers-ms/internal/presentation/http/middlewares"
)

type LoginRoutes struct {
	authHandler       handler.LoginCustomer
	loggingMiddleware *middlewares.Logging
}

func NewLoginRoutes(
	authHandler handler.LoginCustomer,
	loggingMiddleware *middlewares.Logging,
) *LoginRoutes {
	return &LoginRoutes{
		authHandler:       authHandler,
		loggingMiddleware: loggingMiddleware,
	}
}

func (r *LoginRoutes) Routes() helpers.RouteType {
	prefix := "/v1/api"
	return helpers.WithPrefix(prefix, helpers.RouteType{
		{
			Method:      http.MethodPost,
			Path:        "/login",
			Handler:     helpers.TraceHandler(http.MethodPost, prefix+"/login", r.authHandler.Handle),
			Description: "Generate Token",
			Middlewares: helpers.Middlewares{
				r.loggingMiddleware.LoggingMiddleware(),
			},
		},
	})
}
