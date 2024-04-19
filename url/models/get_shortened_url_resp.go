// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// GetShortenedURLResp getShortenedUrlResp
// swagger:model getShortenedUrlResp
type GetShortenedURLResp struct {

	// Shortened URL
	ShortURL string `json:"ShortURL,omitempty"`
}

// Validate validates this get shortened Url resp
func (m *GetShortenedURLResp) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *GetShortenedURLResp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetShortenedURLResp) UnmarshalBinary(b []byte) error {
	var res GetShortenedURLResp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
