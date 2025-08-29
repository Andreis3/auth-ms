package routes

import (
	"net/http"

	"github.com/andreis3/auth-ms/internal/auth/adapter/input/http/helpers"
	"github.com/andreis3/auth-ms/internal/auth/adapter/input/http/middlewares"
	"github.com/andreis3/auth-ms/internal/auth/infra/factory/http/handler"
)

type User struct {
	CreateAuthUser    *handler.CreateAuthUser
	loggingMiddleware *middlewares.Logging
}

func NewUser(
	CreateAuthUser *handler.CreateAuthUser,
	loggingMiddleware *middlewares.Logging,
) *User {
	return &User{
		CreateAuthUser:    CreateAuthUser,
		loggingMiddleware: loggingMiddleware,
	}
}

func (cr *User) Routes() helpers.RouteType {
	prefix := "/auth"
	return helpers.WithPrefix(prefix, helpers.RouteType{
		{
			Method: http.MethodPost,
			Path:   "/signup",
			Handler: helpers.TraceHandler(http.MethodPost, prefix+"/signup", func(w http.ResponseWriter, r *http.Request) {
				cr.CreateAuthUser.NewCreateAuthUser().Handle(w, r)
			}),
			Description: "Create Customer",
			Middlewares: helpers.Middlewares{
				cr.loggingMiddleware.LoggingMiddleware(),
			},
		},
	})
}
