// Code generated by go-swagger; DO NOT EDIT.

package mono_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// BulkInvitationResponse bulk invitation response
//
// swagger:model BulkInvitationResponse
type BulkInvitationResponse struct {

	// failures
	Failures []*BulkInvitationResponseFailuresItems0 `json:"failures"`

	// successes
	Successes []string `json:"successes"`
}

// Validate validates this bulk invitation response
func (m *BulkInvitationResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFailures(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BulkInvitationResponse) validateFailures(formats strfmt.Registry) error {
	if swag.IsZero(m.Failures) { // not required
		return nil
	}

	for i := 0; i < len(m.Failures); i++ {
		if swag.IsZero(m.Failures[i]) { // not required
			continue
		}

		if m.Failures[i] != nil {
			if err := m.Failures[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("failures" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this bulk invitation response based on the context it is used
func (m *BulkInvitationResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFailures(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BulkInvitationResponse) contextValidateFailures(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Failures); i++ {

		if m.Failures[i] != nil {
			if err := m.Failures[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("failures" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *BulkInvitationResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BulkInvitationResponse) UnmarshalBinary(b []byte) error {
	var res BulkInvitationResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// BulkInvitationResponseFailuresItems0 bulk invitation response failures items0
//
// swagger:model BulkInvitationResponseFailuresItems0
type BulkInvitationResponseFailuresItems0 struct {

	// email
	Email string `json:"email,omitempty"`

	// reason
	Reason string `json:"reason,omitempty"`
}

// Validate validates this bulk invitation response failures items0
func (m *BulkInvitationResponseFailuresItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this bulk invitation response failures items0 based on context it is used
func (m *BulkInvitationResponseFailuresItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *BulkInvitationResponseFailuresItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BulkInvitationResponseFailuresItems0) UnmarshalBinary(b []byte) error {
	var res BulkInvitationResponseFailuresItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}