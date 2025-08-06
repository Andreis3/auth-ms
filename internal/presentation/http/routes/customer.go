package routes

import (
	"net/http"

	"github.com/andreis3/customers-ms/internal/infra/factories/presentation/handler-factory"
	"github.com/andreis3/customers-ms/internal/infra/factories/presentation/middleware-factory"
	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
	"github.com/andreis3/customers-ms/internal/presentation/http/middlewares"
)

type CustomerRoutes struct {
	createCustomerHandler *handler_factory.CreateCustomer
	getAddress            *handler_factory.GetCustomerAddresses
	validateCustomer      *middleware_factory.ValidateCustomerFactory
	loggingMiddleware     *middlewares.Logging
}

func NewCustomer(
	createCustomerHandler *handler_factory.CreateCustomer,
	getAddress *handler_factory.GetCustomerAddresses,
	validateCustomer *middleware_factory.ValidateCustomerFactory,
	loggingMiddleware *middlewares.Logging,
) *CustomerRoutes {
	return &CustomerRoutes{
		createCustomerHandler: createCustomerHandler,
		getAddress:            getAddress,
		validateCustomer:      validateCustomer,
		loggingMiddleware:     loggingMiddleware,
	}
}

func (cr *CustomerRoutes) Routes() helpers.RouteType {
	prefix := "/v1/api"
	return helpers.WithPrefix(prefix, helpers.RouteType{
		{
			Method: http.MethodPost,
			Path:   "/customer",
			Handler: helpers.TraceHandler(http.MethodPost, prefix+"/customer", func(w http.ResponseWriter, r *http.Request) {
				cr.createCustomerHandler.NewCreateCustomer().Handle(w, r)
			}),
			Description: "Create Customer",
			Middlewares: helpers.Middlewares{
				cr.loggingMiddleware.LoggingMiddleware(),
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/customer/addresses",
			Handler: helpers.TraceHandler(http.MethodGet, prefix+"/customer/addresses", func(w http.ResponseWriter, r *http.Request) {
				cr.getAddress.NewGetCustomerAddresses().Handle(w, r)
			}),
			Description: "Get Customer Addresses",
			Middlewares: helpers.Middlewares{
				cr.loggingMiddleware.LoggingMiddleware(),
				middleware_factory.MakeValidateCustomerMiddleware(cr.validateCustomer),
			},
		},
	})
}
