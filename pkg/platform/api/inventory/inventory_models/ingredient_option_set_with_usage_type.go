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

// IngredientOptionSetWithUsageType Ingredient Option Set with Usage Type
//
// An ingredient option set paired with how it is used by the ingredient version
//
// swagger:model ingredientOptionSetWithUsageType
type IngredientOptionSetWithUsageType struct {
	IngredientOptionSetWithUsageTypeAllOf0

	IngredientOptionSet
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *IngredientOptionSetWithUsageType) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 IngredientOptionSetWithUsageTypeAllOf0
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.IngredientOptionSetWithUsageTypeAllOf0 = aO0

	// AO1
	var aO1 IngredientOptionSet
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.IngredientOptionSet = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m IngredientOptionSetWithUsageType) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.IngredientOptionSetWithUsageTypeAllOf0)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.IngredientOptionSet)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this ingredient option set with usage type
func (m *IngredientOptionSetWithUsageType) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with IngredientOptionSetWithUsageTypeAllOf0
	if err := m.IngredientOptionSetWithUsageTypeAllOf0.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with IngredientOptionSet
	if err := m.IngredientOptionSet.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this ingredient option set with usage type based on the context it is used
func (m *IngredientOptionSetWithUsageType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with IngredientOptionSetWithUsageTypeAllOf0
	if err := m.IngredientOptionSetWithUsageTypeAllOf0.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with IngredientOptionSet
	if err := m.IngredientOptionSet.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *IngredientOptionSetWithUsageType) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IngredientOptionSetWithUsageType) UnmarshalBinary(b []byte) error {
	var res IngredientOptionSetWithUsageType
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}