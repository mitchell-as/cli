// Code generated by go-swagger; DO NOT EDIT.

package secrets

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	models "github.com/ActiveState/cli/internal/secrets-api/models"
	"github.com/go-openapi/runtime"
	strfmt "github.com/go-openapi/strfmt"
)

// SaveAllUserSecretsReader is a Reader for the SaveAllUserSecrets structure.
type SaveAllUserSecretsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SaveAllUserSecretsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 204:
		result := NewSaveAllUserSecretsNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 401:
		result := NewSaveAllUserSecretsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewSaveAllUserSecretsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewSaveAllUserSecretsNoContent creates a SaveAllUserSecretsNoContent with default headers values
func NewSaveAllUserSecretsNoContent() *SaveAllUserSecretsNoContent {
	return &SaveAllUserSecretsNoContent{}
}

/*SaveAllUserSecretsNoContent handles this case with default header values.

Success
*/
type SaveAllUserSecretsNoContent struct {
}

func (o *SaveAllUserSecretsNoContent) Error() string {
	return fmt.Sprintf("[PATCH /organizations/{organizationID}/user_secrets][%d] saveAllUserSecretsNoContent ", 204)
}

func (o *SaveAllUserSecretsNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewSaveAllUserSecretsUnauthorized creates a SaveAllUserSecretsUnauthorized with default headers values
func NewSaveAllUserSecretsUnauthorized() *SaveAllUserSecretsUnauthorized {
	return &SaveAllUserSecretsUnauthorized{}
}

/*SaveAllUserSecretsUnauthorized handles this case with default header values.

Invalid credentials
*/
type SaveAllUserSecretsUnauthorized struct {
	Payload *models.Message
}

func (o *SaveAllUserSecretsUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /organizations/{organizationID}/user_secrets][%d] saveAllUserSecretsUnauthorized  %+v", 401, o.Payload)
}

func (o *SaveAllUserSecretsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Message)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSaveAllUserSecretsInternalServerError creates a SaveAllUserSecretsInternalServerError with default headers values
func NewSaveAllUserSecretsInternalServerError() *SaveAllUserSecretsInternalServerError {
	return &SaveAllUserSecretsInternalServerError{}
}

/*SaveAllUserSecretsInternalServerError handles this case with default header values.

Server Error
*/
type SaveAllUserSecretsInternalServerError struct {
	Payload *models.Message
}

func (o *SaveAllUserSecretsInternalServerError) Error() string {
	return fmt.Sprintf("[PATCH /organizations/{organizationID}/user_secrets][%d] saveAllUserSecretsInternalServerError  %+v", 500, o.Payload)
}

func (o *SaveAllUserSecretsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Message)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}