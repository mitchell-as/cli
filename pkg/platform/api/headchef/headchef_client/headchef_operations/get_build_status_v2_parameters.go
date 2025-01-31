// Code generated by go-swagger; DO NOT EDIT.

package headchef_operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetBuildStatusV2Params creates a new GetBuildStatusV2Params object
// with the default values initialized.
func NewGetBuildStatusV2Params() *GetBuildStatusV2Params {
	var ()
	return &GetBuildStatusV2Params{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetBuildStatusV2ParamsWithTimeout creates a new GetBuildStatusV2Params object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetBuildStatusV2ParamsWithTimeout(timeout time.Duration) *GetBuildStatusV2Params {
	var ()
	return &GetBuildStatusV2Params{

		timeout: timeout,
	}
}

// NewGetBuildStatusV2ParamsWithContext creates a new GetBuildStatusV2Params object
// with the default values initialized, and the ability to set a context for a request
func NewGetBuildStatusV2ParamsWithContext(ctx context.Context) *GetBuildStatusV2Params {
	var ()
	return &GetBuildStatusV2Params{

		Context: ctx,
	}
}

// NewGetBuildStatusV2ParamsWithHTTPClient creates a new GetBuildStatusV2Params object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetBuildStatusV2ParamsWithHTTPClient(client *http.Client) *GetBuildStatusV2Params {
	var ()
	return &GetBuildStatusV2Params{
		HTTPClient: client,
	}
}

/*GetBuildStatusV2Params contains all the parameters to send to the API endpoint
for the get build status v2 operation typically these are written to a http.Request
*/
type GetBuildStatusV2Params struct {

	/*BuildPlanID
	  The build plan ID of the build.

	*/
	BuildPlanID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get build status v2 params
func (o *GetBuildStatusV2Params) WithTimeout(timeout time.Duration) *GetBuildStatusV2Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get build status v2 params
func (o *GetBuildStatusV2Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get build status v2 params
func (o *GetBuildStatusV2Params) WithContext(ctx context.Context) *GetBuildStatusV2Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get build status v2 params
func (o *GetBuildStatusV2Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get build status v2 params
func (o *GetBuildStatusV2Params) WithHTTPClient(client *http.Client) *GetBuildStatusV2Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get build status v2 params
func (o *GetBuildStatusV2Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBuildPlanID adds the buildPlanID to the get build status v2 params
func (o *GetBuildStatusV2Params) WithBuildPlanID(buildPlanID strfmt.UUID) *GetBuildStatusV2Params {
	o.SetBuildPlanID(buildPlanID)
	return o
}

// SetBuildPlanID adds the buildPlanId to the get build status v2 params
func (o *GetBuildStatusV2Params) SetBuildPlanID(buildPlanID strfmt.UUID) {
	o.BuildPlanID = buildPlanID
}

// WriteToRequest writes these params to a swagger request
func (o *GetBuildStatusV2Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param build_plan_id
	if err := r.SetPathParam("build_plan_id", o.BuildPlanID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
