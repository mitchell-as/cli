// Code generated by go-swagger; DO NOT EDIT.

package buildlogstream_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// TaskSkipped Task Skipped
//
// A message indicating that a requested task for a build has been skipped.
//
// swagger:model taskSkipped
type TaskSkipped struct {

	// reason
	// Required: true
	Reason *string `json:"reason"`

	// task
	// Required: true
	Task *string `json:"task"`
}

// Validate validates this task skipped
func (m *TaskSkipped) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateReason(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTask(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskSkipped) validateReason(formats strfmt.Registry) error {

	if err := validate.Required("reason", "body", m.Reason); err != nil {
		return err
	}

	return nil
}

func (m *TaskSkipped) validateTask(formats strfmt.Registry) error {

	if err := validate.Required("task", "body", m.Task); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TaskSkipped) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskSkipped) UnmarshalBinary(b []byte) error {
	var res TaskSkipped
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
