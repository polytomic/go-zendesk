package zendesk

import (
	"context"
	"encoding/json"
)

type OrganizationFieldAPI interface {
	GetOrganizationFields(ctx context.Context, opts *CustomFieldListOptions) ([]CustomField, Page, error)
}

// GetOrganizationFields returns the list of custom fields for organizations
//
// ref: https://developer.zendesk.com/rest_api/docs/support/organization_fields#list-organization-fields
func (z *Client) GetOrganizationFields(ctx context.Context, opts *CustomFieldListOptions) ([]CustomField, Page, error) {
	var data struct {
		OrganizationFields []CustomField `json:"organization_fields"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &CustomFieldListOptions{}
	}

	u, err := addOptions("/organization_fields.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.OrganizationFields, data.Page, nil
}
