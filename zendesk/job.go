package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// JobStatus provides an "enumeration" for valid Job statuses
type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobWorking   JobStatus = "working"
	JobFailed    JobStatus = "failed"
	JobCompleted JobStatus = "completed"
	JobKilled    JobStatus = "killed"
)

// Terminal returns true if the status is terminal; ie, processing has
// completed.
func (s JobStatus) Terminal() bool {
	switch s {
	case JobFailed, JobCompleted, JobKilled:
		return true
	}
	return false
}

// Job is a Zendesk Job Status
// https://developer.zendesk.com/rest_api/docs/support/job_statuses
type Job struct {
	ID       string                   `json:"id,omitempty"`
	URL      string                   `json:"url,omitempty"`
	Status   JobStatus                `json:"status,omitempty"`
	Total    *int                     `json:"total,omitempty"`
	Progress *int                     `json:"progress,omitempty"`
	Message  *string                  `json:"message,omitempty"`
	Results  []map[string]interface{} `json:"results,omitempty"`
}

// JobStatusAPI is an interface containing the Job status related methods
type JobStatusAPI interface {
	ListJobs(ctx context.Context) ([]Job, error)
	GetJob(ctx context.Context, id string) (Job, error)
	GetJobs(ctx context.Context, ids ...string) ([]Job, error)
}

// ListJobs returns all running Jobs from Zendesk
func (z *Client) ListJobs(ctx context.Context) ([]Job, error) {
	var data struct {
		JobStatuses []Job `json:"job_statuses"`
	}
	u, err := addOptions("/job_statuses.json", struct{}{})
	if err != nil {
		return nil, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.JobStatuses, nil
}

// GetJob returns the status of a background job.
//
// A job may no longer exist to query. Zendesk only logs the last 100 jobs. Jobs
// also expire within an hour.
func (z *Client) GetJob(ctx context.Context, id string) (Job, error) {
	var data struct {
		JobStatus Job `json:"job_status"`
	}
	u, err := addOptions(fmt.Sprintf("/job_statuses/%s.json", id), struct{}{})
	if err != nil {
		return Job{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return Job{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return Job{}, err
	}
	return data.JobStatus, nil
}

// GetJobs returns the status of one or more background jobs.
func (z *Client) GetJobs(ctx context.Context, ids ...string) ([]Job, error) {
	opts := struct {
		IDs string `json:"ids"`
	}{
		IDs: strings.Join(ids, ","),
	}
	var data struct {
		JobStatuses []Job `json:"job_statuses"`
	}
	u, err := addOptions("/job_statuses/show_many.json", opts)
	if err != nil {
		return nil, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.JobStatuses, nil
}
