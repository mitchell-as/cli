// Code generated by go-swagger; DO NOT EDIT.

package inventory_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// IngredientOptionSetRevisionCoreAllOf0 ingredient option set revision core all of0
//
// swagger:model ingredientOptionSetRevisionCoreAllOf0
type IngredientOptionSetRevisionCoreAllOf0 struct {

	// The ingredient options. May be an empty list if this set is intended to override and disable another set in the same group.
	// Required: true
	IngredientOptions []*IngredientOption `json:"ingredient_options"`
}

// Validate validates this ingredient option set revision core all of0
func (m *IngredientOptionSetRevisionCoreAllOf0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateIngredientOptions(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IngredientOptionSetRevisionCoreAllOf0) validateIngredientOptions(formats strfmt.Registry) error {

	if err := validate.Required("ingredient_options", "body", m.IngredientOptions); err != nil {
		return err
	}

	for i := 0; i < len(m.IngredientOptions); i++ {
		if swag.IsZero(m.IngredientOptions[i]) { // not required
			continue
		}

		if m.IngredientOptions[i] != nil {
			if err := m.IngredientOptions[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ingredient_options" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *IngredientOptionSetRevisionCoreAllOf0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IngredientOptionSetRevisionCoreAllOf0) UnmarshalBinary(b []byte) error {
	var res IngredientOptionSetRevisionCoreAllOf0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}