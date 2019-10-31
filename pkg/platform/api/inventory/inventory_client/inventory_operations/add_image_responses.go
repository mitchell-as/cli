// Code generated by go-swagger; DO NOT EDIT.

package inventory_operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	inventory_models "github.com/ActiveState/cli/pkg/platform/api/inventory/inventory_models"
)

// AddImageReader is a Reader for the AddImage structure.
type AddImageReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddImageReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewAddImageCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewAddImageBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewAddImageDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAddImageCreated creates a AddImageCreated with default headers values
func NewAddImageCreated() *AddImageCreated {
	return &AddImageCreated{}
}

/*AddImageCreated handles this case with default header values.

The added image
*/
type AddImageCreated struct {
	Payload *inventory_models.V1Image
}

func (o *AddImageCreated) Error() string {
	return fmt.Sprintf("[POST /v1/images][%d] addImageCreated  %+v", 201, o.Payload)
}

func (o *AddImageCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(inventory_models.V1Image)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddImageBadRequest creates a AddImageBadRequest with default headers values
func NewAddImageBadRequest() *AddImageBadRequest {
	return &AddImageBadRequest{}
}

/*AddImageBadRequest handles this case with default header values.

If the image is invalid
*/
type AddImageBadRequest struct {
	Payload *inventory_models.RestAPIValidationError
}

func (o *AddImageBadRequest) Error() string {
	return fmt.Sprintf("[POST /v1/images][%d] addImageBadRequest  %+v", 400, o.Payload)
}

func (o *AddImageBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(inventory_models.RestAPIValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddImageDefault creates a AddImageDefault with default headers values
func NewAddImageDefault(code int) *AddImageDefault {
	return &AddImageDefault{
		_statusCode: code,
	}
}

/*AddImageDefault handles this case with default header values.

If there is an error processing the image
*/
type AddImageDefault struct {
	_statusCode int

	Payload *inventory_models.RestAPIError
}

// Code gets the status code for the add image default response
func (o *AddImageDefault) Code() int {
	return o._statusCode
}

func (o *AddImageDefault) Error() string {
	return fmt.Sprintf("[POST /v1/images][%d] addImage default  %+v", o._statusCode, o.Payload)
}

func (o *AddImageDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(inventory_models.RestAPIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}