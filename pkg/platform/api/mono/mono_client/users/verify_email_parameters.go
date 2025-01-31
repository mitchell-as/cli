// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewVerifyEmailParams creates a new VerifyEmailParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewVerifyEmailParams() *VerifyEmailParams {
	return &VerifyEmailParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewVerifyEmailParamsWithTimeout creates a new VerifyEmailParams object
// with the ability to set a timeout on a request.
func NewVerifyEmailParamsWithTimeout(timeout time.Duration) *VerifyEmailParams {
	return &VerifyEmailParams{
		timeout: timeout,
	}
}

// NewVerifyEmailParamsWithContext creates a new VerifyEmailParams object
// with the ability to set a context for a request.
func NewVerifyEmailParamsWithContext(ctx context.Context) *VerifyEmailParams {
	return &VerifyEmailParams{
		Context: ctx,
	}
}

// NewVerifyEmailParamsWithHTTPClient creates a new VerifyEmailParams object
// with the ability to set a custom HTTPClient for a request.
func NewVerifyEmailParamsWithHTTPClient(client *http.Client) *VerifyEmailParams {
	return &VerifyEmailParams{
		HTTPClient: client,
	}
}

/* VerifyEmailParams contains all the parameters to send to the API endpoint
   for the verify email operation.

   Typically these are written to a http.Request.
*/
type VerifyEmailParams struct {

	/* Email.

	   email address to change
	*/
	Email string

	/* Username.

	   username of desired User
	*/
	Username string

	/* VerifyCode.

	   verification code for email
	*/
	VerifyCode string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the verify email params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *VerifyEmailParams) WithDefaults() *VerifyEmailParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the verify email params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *VerifyEmailParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the verify email params
func (o *VerifyEmailParams) WithTimeout(timeout time.Duration) *VerifyEmailParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the verify email params
func (o *VerifyEmailParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the verify email params
func (o *VerifyEmailParams) WithContext(ctx context.Context) *VerifyEmailParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the verify email params
func (o *VerifyEmailParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the verify email params
func (o *VerifyEmailParams) WithHTTPClient(client *http.Client) *VerifyEmailParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the verify email params
func (o *VerifyEmailParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEmail adds the email to the verify email params
func (o *VerifyEmailParams) WithEmail(email string) *VerifyEmailParams {
	o.SetEmail(email)
	return o
}

// SetEmail adds the email to the verify email params
func (o *VerifyEmailParams) SetEmail(email string) {
	o.Email = email
}

// WithUsername adds the username to the verify email params
func (o *VerifyEmailParams) WithUsername(username string) *VerifyEmailParams {
	o.SetUsername(username)
	return o
}

// SetUsername adds the username to the verify email params
func (o *VerifyEmailParams) SetUsername(username string) {
	o.Username = username
}

// WithVerifyCode adds the verifyCode to the verify email params
func (o *VerifyEmailParams) WithVerifyCode(verifyCode string) *VerifyEmailParams {
	o.SetVerifyCode(verifyCode)
	return o
}

// SetVerifyCode adds the verifyCode to the verify email params
func (o *VerifyEmailParams) SetVerifyCode(verifyCode string) {
	o.VerifyCode = verifyCode
}

// WriteToRequest writes these params to a swagger request
func (o *VerifyEmailParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param email
	if err := r.SetPathParam("email", o.Email); err != nil {
		return err
	}

	// path param username
	if err := r.SetPathParam("username", o.Username); err != nil {
		return err
	}

	// query param verifyCode
	qrVerifyCode := o.VerifyCode
	qVerifyCode := qrVerifyCode
	if qVerifyCode != "" {

		if err := r.SetQueryParam("verifyCode", qVerifyCode); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
