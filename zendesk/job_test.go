package zendesk

import (
	"net/http"
	"testing"
)

func TestListJobs(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodGet, "list_jobs.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	jobs, err := client.ListJobs(ctx)
	if err != nil {
		t.Fatalf("Failed to list jobs: %s", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected length of jobs is 2, but got %d", len(jobs))
	}
}

func TestGetJob(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodGet, "job.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	job, err := client.GetJob(ctx, "8b726e606741012ffc2d782bcb7848fe")
	if err != nil {
		t.Fatalf("Failed to get job: %s", err)
	}

	expectedID := "e57bc8b851f93833d87748557c9b7e4b"
	if job.ID != expectedID {
		t.Fatalf("Returned job does not have the expected ID %s. Job id is %s", expectedID, job.ID)
	}
	if job.Message == nil || *job.Message == "" {
		t.Fatalf(("Returned job does not have message populated"))
	}
}

func TestGetJobsById(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodGet, "jobs.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	jobs, err := client.GetJobs(ctx, "8b726e606741012ffc2d782bcb7848fe")
	if err != nil {
		t.Fatalf("Failed to get jobs: %s", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected length of jobs is 1, but got %d", len(jobs))
	}

	expectedID := "8b726e606741012ffc2d782bcb7848fe"
	if jobs[0].ID != expectedID {
		t.Fatalf("Returned job does not have the expected ID %s. User id is %s", expectedID, jobs[0].ID)
	}
}
