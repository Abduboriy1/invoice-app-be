// internal/interfaces/http/dto/jira.go
package dto

type PullWorklogsRequest struct {
	StartDate string   `json:"start_date" validate:"required"` // YYYY-MM-DD
	EndDate   string   `json:"end_date" validate:"required"`   // YYYY-MM-DD
	IssueKeys []string `json:"issue_keys"`                     // Optional: specific issues
}

type PullIssueWorklogsRequest struct {
	IssueKey string `json:"issue_key" validate:"required"`
}

type PushWorklogRequest struct {
	TimeEntryID string `json:"time_entry_id" validate:"required"`
	IssueKey    string `json:"issue_key" validate:"required"`
}

type JiraConfigRequest struct {
	BaseURL  string `json:"base_url" validate:"required,url"`
	Email    string `json:"email" validate:"required,email"`
	APIToken string `json:"api_token" validate:"required"`
}
