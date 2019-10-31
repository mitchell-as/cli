// Code generated by go-swagger; DO NOT EDIT.

package inventory_operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetIngredientVersionsParams creates a new GetIngredientVersionsParams object
// with the default values initialized.
func NewGetIngredientVersionsParams() *GetIngredientVersionsParams {
	var (
		allowUnstableDefault = bool(false)
		limitDefault         = int64(50)
		pageDefault          = int64(1)
	)
	return &GetIngredientVersionsParams{
		AllowUnstable: &allowUnstableDefault,
		Limit:         &limitDefault,
		Page:          &pageDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewGetIngredientVersionsParamsWithTimeout creates a new GetIngredientVersionsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetIngredientVersionsParamsWithTimeout(timeout time.Duration) *GetIngredientVersionsParams {
	var (
		allowUnstableDefault = bool(false)
		limitDefault         = int64(50)
		pageDefault          = int64(1)
	)
	return &GetIngredientVersionsParams{
		AllowUnstable: &allowUnstableDefault,
		Limit:         &limitDefault,
		Page:          &pageDefault,

		timeout: timeout,
	}
}

// NewGetIngredientVersionsParamsWithContext creates a new GetIngredientVersionsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetIngredientVersionsParamsWithContext(ctx context.Context) *GetIngredientVersionsParams {
	var (
		allowUnstableDefault = bool(false)
		limitDefault         = int64(50)
		pageDefault          = int64(1)
	)
	return &GetIngredientVersionsParams{
		AllowUnstable: &allowUnstableDefault,
		Limit:         &limitDefault,
		Page:          &pageDefault,

		Context: ctx,
	}
}

// NewGetIngredientVersionsParamsWithHTTPClient creates a new GetIngredientVersionsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetIngredientVersionsParamsWithHTTPClient(client *http.Client) *GetIngredientVersionsParams {
	var (
		allowUnstableDefault = bool(false)
		limitDefault         = int64(50)
		pageDefault          = int64(1)
	)
	return &GetIngredientVersionsParams{
		AllowUnstable: &allowUnstableDefault,
		Limit:         &limitDefault,
		Page:          &pageDefault,
		HTTPClient:    client,
	}
}

/*GetIngredientVersionsParams contains all the parameters to send to the API endpoint
for the get ingredient versions operation typically these are written to a http.Request
*/
type GetIngredientVersionsParams struct {

	/*AllowUnstable
	  Whether to show an unstable revision of a resource if there is an available unstable version newer than the newest available stable version

	*/
	AllowUnstable *bool
	/*IngredientID*/
	IngredientID strfmt.UUID
	/*Limit
	  The maximum number of items returned per page

	*/
	Limit *int64
	/*Page
	  The page number returned

	*/
	Page *int64
	/*StateAt
	  Show the state of a resource as it was at the specified timestamp. If omitted, shows the current state of the resource.

	*/
	StateAt *strfmt.DateTime

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithTimeout(timeout time.Duration) *GetIngredientVersionsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithContext(ctx context.Context) *GetIngredientVersionsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithHTTPClient(client *http.Client) *GetIngredientVersionsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAllowUnstable adds the allowUnstable to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithAllowUnstable(allowUnstable *bool) *GetIngredientVersionsParams {
	o.SetAllowUnstable(allowUnstable)
	return o
}

// SetAllowUnstable adds the allowUnstable to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetAllowUnstable(allowUnstable *bool) {
	o.AllowUnstable = allowUnstable
}

// WithIngredientID adds the ingredientID to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithIngredientID(ingredientID strfmt.UUID) *GetIngredientVersionsParams {
	o.SetIngredientID(ingredientID)
	return o
}

// SetIngredientID adds the ingredientId to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetIngredientID(ingredientID strfmt.UUID) {
	o.IngredientID = ingredientID
}

// WithLimit adds the limit to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithLimit(limit *int64) *GetIngredientVersionsParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithPage adds the page to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithPage(page *int64) *GetIngredientVersionsParams {
	o.SetPage(page)
	return o
}

// SetPage adds the page to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetPage(page *int64) {
	o.Page = page
}

// WithStateAt adds the stateAt to the get ingredient versions params
func (o *GetIngredientVersionsParams) WithStateAt(stateAt *strfmt.DateTime) *GetIngredientVersionsParams {
	o.SetStateAt(stateAt)
	return o
}

// SetStateAt adds the stateAt to the get ingredient versions params
func (o *GetIngredientVersionsParams) SetStateAt(stateAt *strfmt.DateTime) {
	o.StateAt = stateAt
}

// WriteToRequest writes these params to a swagger request
func (o *GetIngredientVersionsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.AllowUnstable != nil {

		// query param allow_unstable
		var qrAllowUnstable bool
		if o.AllowUnstable != nil {
			qrAllowUnstable = *o.AllowUnstable
		}
		qAllowUnstable := swag.FormatBool(qrAllowUnstable)
		if qAllowUnstable != "" {
			if err := r.SetQueryParam("allow_unstable", qAllowUnstable); err != nil {
				return err
			}
		}

	}

	// path param ingredient_id
	if err := r.SetPathParam("ingredient_id", o.IngredientID.String()); err != nil {
		return err
	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int64
		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {
			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}

	}

	if o.Page != nil {

		// query param page
		var qrPage int64
		if o.Page != nil {
			qrPage = *o.Page
		}
		qPage := swag.FormatInt64(qrPage)
		if qPage != "" {
			if err := r.SetQueryParam("page", qPage); err != nil {
				return err
			}
		}

	}

	if o.StateAt != nil {

		// query param state_at
		var qrStateAt strfmt.DateTime
		if o.StateAt != nil {
			qrStateAt = *o.StateAt
		}
		qStateAt := qrStateAt.String()
		if qStateAt != "" {
			if err := r.SetQueryParam("state_at", qStateAt); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}