// Code generated by go-swagger; DO NOT EDIT.

package endpoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PutEndpointIDLabelsHandlerFunc turns a function with the right signature into a put endpoint ID labels handler
type PutEndpointIDLabelsHandlerFunc func(PutEndpointIDLabelsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutEndpointIDLabelsHandlerFunc) Handle(params PutEndpointIDLabelsParams) middleware.Responder {
	return fn(params)
}

// PutEndpointIDLabelsHandler interface for that can handle valid put endpoint ID labels params
type PutEndpointIDLabelsHandler interface {
	Handle(PutEndpointIDLabelsParams) middleware.Responder
}

// NewPutEndpointIDLabels creates a new http.Handler for the put endpoint ID labels operation
func NewPutEndpointIDLabels(ctx *middleware.Context, handler PutEndpointIDLabelsHandler) *PutEndpointIDLabels {
	return &PutEndpointIDLabels{Context: ctx, Handler: handler}
}

/*PutEndpointIDLabels swagger:route PUT /endpoint/{id}/labels endpoint putEndpointIdLabels

Modify label configuration of endpoint

Updates the list of labels associated with an endpoint by applying
a label modificator structure to the label configuration of an
endpoint.

The label configuration mutation is only executed as a whole, i.e.
if any of the labels to be deleted are not either on the list of
orchestration system labels, custom labels, or already disabled,
then the request will fail. Labels to be added which already exist
on either the orchestration list or custom list will be ignored.


*/
type PutEndpointIDLabels struct {
	Context *middleware.Context
	Handler PutEndpointIDLabelsHandler
}

func (o *PutEndpointIDLabels) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutEndpointIDLabelsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
