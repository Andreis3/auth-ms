package routes

import (
	"net/http"

	helpers2 "github.com/andreis3/auth-ms/internal/adapter/input/http/helpers"
	"github.com/andreis3/auth-ms/internal/adapter/input/http/middlewares"
	"github.com/andreis3/auth-ms/internal/infra/factory/http/handler"
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

func (cr *User) Routes() helpers2.RouteType {
	prefix := "/auth"
	return helpers2.WithPrefix(prefix, helpers2.RouteType{
		{
			Method: http.MethodPost,
			Path:   "/signup",
			Handler: helpers2.TraceHandler(http.MethodPost, prefix+"/signup", func(w http.ResponseWriter, r *http.Request) {
				cr.CreateAuthUser.NewCreateAuthUser().Handle(w, r)
			}),
			Description: "Create Customer",
			Middlewares: helpers2.Middlewares{
				cr.loggingMiddleware.LoggingMiddleware(),
			},
		},
	})
}
