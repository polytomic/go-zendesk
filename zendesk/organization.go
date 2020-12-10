package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Organization is struct for organization payload
// https://developer.zendesk.com/rest_api/docs/support/organizations
type Organization struct {
	ID                 int64                  `json:"id,omitempty"`
	ExternalID         string                 `json:"external_id,omitempty"`
	URL                string                 `json:"url,omitempty"`
	Name               string                 `json:"name"`
	DomainNames        []string               `json:"domain_names"`
	GroupID            int64                  `json:"group_id"`
	SharedTickets      bool                   `json:"shared_tickets"`
	SharedComments     bool                   `json:"shared_comments"`
	Tags               []string               `json:"tags"`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
	UpdatedAt          time.Time              `json:"updated_at,omitempty"`
	OrganizationFields map[string]interface{} `json:"organization_fields,omitempty"`
}

// OrganizationAPI an interface containing all methods associated with zendesk organizations
type OrganizationAPI interface {
	CreateOrganization(ctx context.Context, org Organization) (Organization, error)
	GetOrganization(ctx context.Context, orgID int64) (Organization, error)
	GetOrganizations(ctx context.Context, orgIDs ...int64) ([]Organization, error)
	GetOrganizationsByExternalID(ctx context.Context, externalOrgIds ...int64) ([]Organization, error)
	UpdateOrganization(ctx context.Context, orgID int64, org Organization) (Organization, error)
	UpdateManyOrganizations(ctx context.Context, organizations []Organization) (Job, error)
	DeleteOrganization(ctx context.Context, orgID int64) error
}

// CreateOrganization creates new organization
// https://developer.zendesk.com/rest_api/docs/support/organizations#create-organization
func (z *Client) CreateOrganization(ctx context.Context, org Organization) (Organization, error) {
	var data, result struct {
		Organization Organization `json:"organization"`
	}

	data.Organization = org

	body, err := z.post(ctx, "/organizations.json", data)
	if err != nil {
		return Organization{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Organization{}, err
	}

	return result.Organization, nil
}

// GetOrganization gets a specified organization
// ref: https://developer.zendesk.com/rest_api/docs/support/organizations#show-organization
func (z *Client) GetOrganization(ctx context.Context, orgID int64) (Organization, error) {
	var result struct {
		Organization Organization `json:"organization"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/organizations/%d.json", orgID))

	if err != nil {
		return Organization{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Organization{}, err
	}

	return result.Organization, err
}

// GetOrganizations retrieves one or more organizations by their IDs
// ref: https://developer.zendesk.com/rest_api/docs/support/organizations#show-many-organizations
func (z *Client) GetOrganizations(ctx context.Context, orgIDs ...int64) ([]Organization, error) {
	var result struct {
		Organizations []Organization `json:"organizations"`
	}

	ids := make([]string, len(orgIDs))
	for i, id := range orgIDs {
		ids[i] = strconv.FormatInt(id, 10)
	}
	query := url.Values{}
	query.Add("ids", strings.Join(ids, ","))
	body, err := z.get(ctx, fmt.Sprintf("/organizations/show_many.json?%s", query.Encode()))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Organizations, err
}

// GetOrganizationsByExternalID retrieves one or more organization by their External IDs
// ref: https://developer.zendesk.com/rest_api/docs/support/organizations#show-many-organizations
func (z *Client) GetOrganizationsByExternalID(ctx context.Context, externalOrgIds ...int64) ([]Organization, error) {
	var result struct {
		Organizations []Organization `json:"organizations"`
	}

	ids := make([]string, len(externalOrgIds))
	for i, id := range externalOrgIds {
		ids[i] = strconv.FormatInt(id, 10)
	}
	query := url.Values{}
	query.Add("external_ids", strings.Join(ids, ","))
	body, err := z.get(ctx, fmt.Sprintf("/organizations/show_many.json?%s", query.Encode()))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Organizations, err
}

// UpdateOrganization updates a organization with the specified organization
// ref: https://developer.zendesk.com/rest_api/docs/support/organizations#update-organization
func (z *Client) UpdateOrganization(ctx context.Context, orgID int64, org Organization) (Organization, error) {
	var result, data struct {
		Organization Organization `json:"organization"`
	}

	data.Organization = org

	body, err := z.put(ctx, fmt.Sprintf("/organizations/%d.json", orgID), data)

	if err != nil {
		return Organization{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Organization{}, err
	}

	return result.Organization, err
}

// CreateOrUpdateOrganization either updates an existing organization or creates
// a new one. It returns the organization and a boolean flag that is true if the
// org is newly created.
//
// ref: https://developer.zendesk.com/rest_api/docs/support/organizations#create-or-update-organization
func (z *Client) CreateOrUpdateOrganization(ctx context.Context, org Organization) (Organization, bool, error) {
	var data, result struct {
		Organization Organization `json:"organization"`
	}

	data.Organization = org

	body, resp, err := z.postWithResponse(ctx, "/organizations/create_or_update.json", data)
	if err != nil {
		return Organization{}, false, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return Organization{}, false, Error{
			body: body,
			resp: resp,
		}
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Organization{}, false, err
	}

	return result.Organization, resp.StatusCode == http.StatusCreated, nil
}

// DeleteOrganization deletes the specified organization
// ref: https://developer.zendesk.com/rest_api/docs/support/organizations#delete-organization
func (z *Client) DeleteOrganization(ctx context.Context, orgID int64) error {
	err := z.delete(ctx, fmt.Sprintf("/organizations/%d.json", orgID))

	if err != nil {
		return err
	}

	return nil
}

// UpdateManyOrganizations updates up to 100 organizations with a single request
// via a background job.
//
// ref: https://developer.zendesk.com/rest_api/docs/support/organizations#update-many-organizations
func (z *Client) UpdateManyOrganizations(ctx context.Context, organizations []Organization) (Job, error) {
	data := struct {
		Organizations []Organization `json:"organizations"`
	}{organizations}
	var result struct {
		Job Job `json:"job_status"`
	}

	body, err := z.post(ctx, "/organizations/update_many.json", data, expectStatus(http.StatusOK))
	if err != nil {
		return Job{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Job{}, err
	}
	return result.Job, nil
}
