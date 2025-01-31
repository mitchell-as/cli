// Code generated by go-swagger; DO NOT EDIT.

package inventory_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GpuArchitectureCore GPU Architecture Core
//
// The properties of a GPU architecture needed to create a new one
//
// swagger:model gpuArchitectureCore
type GpuArchitectureCore struct {
	GpuArchitectureCoreAllOf0

	RevisionedFeatureProvider
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *GpuArchitectureCore) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 GpuArchitectureCoreAllOf0
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.GpuArchitectureCoreAllOf0 = aO0

	// AO1
	var aO1 RevisionedFeatureProvider
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.RevisionedFeatureProvider = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m GpuArchitectureCore) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.GpuArchitectureCoreAllOf0)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.RevisionedFeatureProvider)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this gpu architecture core
func (m *GpuArchitectureCore) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with GpuArchitectureCoreAllOf0
	if err := m.GpuArchitectureCoreAllOf0.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with RevisionedFeatureProvider
	if err := m.RevisionedFeatureProvider.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this gpu architecture core based on the context it is used
func (m *GpuArchitectureCore) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with GpuArchitectureCoreAllOf0
	if err := m.GpuArchitectureCoreAllOf0.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with RevisionedFeatureProvider
	if err := m.RevisionedFeatureProvider.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *GpuArchitectureCore) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GpuArchitectureCore) UnmarshalBinary(b []byte) error {
	var res GpuArchitectureCore
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
