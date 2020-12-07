// Code generated by go-swagger; DO NOT EDIT.

package inventory_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ImageAllOf0 image all of0
//
// swagger:model imageAllOf0
type ImageAllOf0 struct {

	// image id
	// Required: true
	// Format: uuid
	ImageID *strfmt.UUID `json:"image_id"`

	// links
	// Required: true
	Links *SelfLink `json:"links"`
}

// Validate validates this image all of0
func (m *ImageAllOf0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateImageID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLinks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ImageAllOf0) validateImageID(formats strfmt.Registry) error {

	if err := validate.Required("image_id", "body", m.ImageID); err != nil {
		return err
	}

	if err := validate.FormatOf("image_id", "body", "uuid", m.ImageID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ImageAllOf0) validateLinks(formats strfmt.Registry) error {

	if err := validate.Required("links", "body", m.Links); err != nil {
		return err
	}

	if m.Links != nil {
		if err := m.Links.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("links")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ImageAllOf0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ImageAllOf0) UnmarshalBinary(b []byte) error {
	var res ImageAllOf0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}