package zendesk

import "time"

// CustomField describes custom fields for users and organizations
type CustomField struct {
	ID                  int64               `json:"id,omitempty"`
	URL                 string              `json:"url,omitempty"`
	Key                 string              `json:"key,omitempty"`
	Type                string              `json:"type"`
	Title               string              `json:"title"`
	RawTitle            string              `json:"raw_title,omitempty"`
	Description         string              `json:"description,omitempty"`
	RawDescription      string              `json:"raw_description,omitempty"`
	Position            int64               `json:"position,omitempty"`
	Active              bool                `json:"active,omitempty"`
	System              bool                `json:"system,omitempty"`
	RegexpForValidation string              `json:"regexp_for_validation,omitempty"`
	Tag                 string              `json:"tag,omitempty"`
	CustomFieldOptions  []CustomFieldOption `json:"custom_field_options"`
	CreatedAt           time.Time           `json:"created_at,omitempty"`
	UpdatedAt           time.Time           `json:"updated_at,omitempty"`
}

// CustomFieldListOptions provides pagination options for custom field lists
type CustomFieldListOptions struct {
	PageOptions
}

// CustomFieldOption is struct for value of `custom_field_options`
type CustomFieldOption struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Position int64  `json:"position,omitempty"`
	RawName  string `json:"raw_name,omitempty"`
	URL      string `json:"url,omitempty"`
	Value    string `json:"value"`
}
