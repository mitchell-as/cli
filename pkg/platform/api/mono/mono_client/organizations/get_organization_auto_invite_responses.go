// Code generated by go-swagger; DO NOT EDIT.

package organizations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/ActiveState/cli/pkg/platform/api/mono/mono_models"
)

// GetOrganizationAutoInviteReader is a Reader for the GetOrganizationAutoInvite structure.
type GetOrganizationAutoInviteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetOrganizationAutoInviteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetOrganizationAutoInviteOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 403:
		result := NewGetOrganizationAutoInviteForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetOrganizationAutoInviteInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetOrganizationAutoInviteOK creates a GetOrganizationAutoInviteOK with default headers values
func NewGetOrganizationAutoInviteOK() *GetOrganizationAutoInviteOK {
	return &GetOrganizationAutoInviteOK{}
}

/* GetOrganizationAutoInviteOK describes a response with status code 200, with default header values.

Success
*/
type GetOrganizationAutoInviteOK struct {
	Payload []*mono_models.AutoInviteSettings
}

func (o *GetOrganizationAutoInviteOK) Error() string {
	return fmt.Sprintf("[GET /organizations/{organizationName}/autoinvite][%d] getOrganizationAutoInviteOK  %+v", 200, o.Payload)
}
func (o *GetOrganizationAutoInviteOK) GetPayload() []*mono_models.AutoInviteSettings {
	return o.Payload
}

func (o *GetOrganizationAutoInviteOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetOrganizationAutoInviteForbidden creates a GetOrganizationAutoInviteForbidden with default headers values
func NewGetOrganizationAutoInviteForbidden() *GetOrganizationAutoInviteForbidden {
	return &GetOrganizationAutoInviteForbidden{}
}

/* GetOrganizationAutoInviteForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type GetOrganizationAutoInviteForbidden struct {
	Payload *mono_models.Message
}

func (o *GetOrganizationAutoInviteForbidden) Error() string {
	return fmt.Sprintf("[GET /organizations/{organizationName}/autoinvite][%d] getOrganizationAutoInviteForbidden  %+v", 403, o.Payload)
}
func (o *GetOrganizationAutoInviteForbidden) GetPayload() *mono_models.Message {
	return o.Payload
}

func (o *GetOrganizationAutoInviteForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(mono_models.Message)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetOrganizationAutoInviteInternalServerError creates a GetOrganizationAutoInviteInternalServerError with default headers values
func NewGetOrganizationAutoInviteInternalServerError() *GetOrganizationAutoInviteInternalServerError {
	return &GetOrganizationAutoInviteInternalServerError{}
}

/* GetOrganizationAutoInviteInternalServerError describes a response with status code 500, with default header values.

Server Error
*/
type GetOrganizationAutoInviteInternalServerError struct {
	Payload *mono_models.Message
}

func (o *GetOrganizationAutoInviteInternalServerError) Error() string {
	return fmt.Sprintf("[GET /organizations/{organizationName}/autoinvite][%d] getOrganizationAutoInviteInternalServerError  %+v", 500, o.Payload)
}
func (o *GetOrganizationAutoInviteInternalServerError) GetPayload() *mono_models.Message {
	return o.Payload
}

func (o *GetOrganizationAutoInviteInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(mono_models.Message)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}