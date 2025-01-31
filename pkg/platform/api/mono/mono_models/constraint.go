// Code generated by go-swagger; DO NOT EDIT.

package mono_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Constraint constraint
//
// swagger:model Constraint
type Constraint struct {

	// comparator
	Comparator string `json:"comparator,omitempty"`

	// version
	Version string `json:"version,omitempty"`
}

// Validate validates this constraint
func (m *Constraint) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this constraint based on context it is used
func (m *Constraint) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Constraint) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Constraint) UnmarshalBinary(b []byte) error {
	var res Constraint
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
