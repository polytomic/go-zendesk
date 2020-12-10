package zendesk

import (
	"context"
	"encoding/json"
)

type UserFieldAPI interface {
	GetUserFields(ctx context.Context, opts *CustomFieldListOptions) ([]CustomField, Page, error)
}

// GetUserFields returns the list of custom fields for Users
//
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#list-user-fields
func (z *Client) GetUserFields(ctx context.Context, opts *CustomFieldListOptions) ([]CustomField, Page, error) {
	var data struct {
		UserFields []CustomField `json:"user_fields"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &CustomFieldListOptions{}
	}

	u, err := addOptions("/user_fields.json", tmp)
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
	return data.UserFields, data.Page, nil
}
