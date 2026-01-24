// internal/infrastructure/integrations/jira/client.go
package jira

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	baseURL string
	email   string
	apiKey  string
	client  *resty.Client
}

func NewClient(baseURL, email, apiKey string) *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetBasicAuth(email, apiKey).
		SetTimeout(30 * time.Second)

	return &Client{
		baseURL: baseURL,
		email:   email,
		apiKey:  apiKey,
		client:  client,
	}
}

type Worklog struct {
	ID               string  `json:"id"`
	IssueKey         string  `json:"-"` // Not from API, set manually
	Author           Author  `json:"author"`
	UpdateAuthor     Author  `json:"updateAuthor"`
	Comment          Comment `json:"comment"`
	Created          Time    `json:"created"`
	Updated          Time    `json:"updated"`
	Started          Time    `json:"started"`
	TimeSpent        string  `json:"timeSpent"`
	TimeSpentSeconds int     `json:"timeSpentSeconds"`
}

type Time struct {
	time.Time
}

// UnmarshalJSON handles Jira's timestamp format
func (t *Time) UnmarshalJSON(b []byte) error {
	// Remove quotes from JSON string
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("invalid time string: %s", s)
	}
	s = s[1 : len(s)-1]

	// Try multiple formats that Jira might use
	formats := []string{
		"2006-01-02T15:04:05.000-0700", // Jira format with milliseconds
		"2006-01-02T15:04:05.000Z0700", // Alternative format
		time.RFC3339,                   // Standard RFC3339
		time.RFC3339Nano,               // RFC3339 with nanoseconds
	}

	var err error
	for _, format := range formats {
		t.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("unable to parse time %q: %w", s, err)
}

type Comment struct {
	Type    string    `json:"type"`
	Version int       `json:"version"`
	Content []Content `json:"content"`
}

type Content struct {
	Type    string        `json:"type"`
	Content []TextContent `json:"content"`
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Author struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// GetWorklogs gets all worklogs for a specific issue
func (c *Client) GetWorklogs(ctx context.Context, issueKey string) ([]Worklog, error) {
	var result struct {
		Worklogs []Worklog `json:"worklogs"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueKey))

	if err != nil {
		return nil, fmt.Errorf("fetching worklogs: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("jira API error: %s - %s", resp.Status(), string(resp.Body()))
	}

	// Set issue key for each worklog
	for i := range result.Worklogs {
		result.Worklogs[i].IssueKey = issueKey
	}

	return result.Worklogs, nil
}

// GetWorklogsByDateRange gets worklogs for multiple issues within a date range
func (c *Client) GetWorklogsByDateRange(ctx context.Context, issueKeys []string, startDate, endDate time.Time) ([]Worklog, error) {
	var allWorklogs []Worklog

	if len(issueKeys) == 0 {
		return allWorklogs, nil
	}

	for _, issueKey := range issueKeys {
		worklogs, err := c.GetWorklogs(ctx, issueKey)
		if err != nil {
			// Log error but continue with other issues
			fmt.Printf("Error fetching worklogs for %s: %v\n", issueKey, err)
			continue
		}

		// Filter by date range
		for _, wl := range worklogs {
			if (wl.Started.Equal(startDate) || wl.Started.After(startDate)) &&
				(wl.Started.Before(endDate.Add(24 * time.Hour))) {
				allWorklogs = append(allWorklogs, wl)
			}
		}
	}

	return allWorklogs, nil
}

// GetIssuesWithWorklogsByDateRange uses the new JQL API to find issues with worklogs in date range
func (c *Client) GetIssuesWithWorklogsByDateRange(ctx context.Context, startDate, endDate time.Time, maxResults int) ([]string, error) {
	jql := fmt.Sprintf("worklogDate >= '%s' AND worklogDate <= '%s' AND assignee = currentUser() ORDER BY updated DESC",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	if maxResults <= 0 {
		maxResults = 100 // Default
	}

	payload := map[string]interface{}{
		"jql":        jql,
		"fields":     []string{"key"},
		"maxResults": maxResults,
		"startAt":    0,
	}

	var result struct {
		Issues []struct {
			Key string `json:"key"`
		} `json:"issues"`
		Total      int `json:"total"`
		MaxResults int `json:"maxResults"`
		StartAt    int `json:"startAt"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&result).
		Post("/rest/api/3/search/jql")

	if err != nil {
		return nil, fmt.Errorf("searching issues: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("jira API error: %s - %s", resp.Status(), string(resp.Body()))
	}

	keys := make([]string, len(result.Issues))
	for i, issue := range result.Issues {
		keys[i] = issue.Key
	}

	return keys, nil
}

// GetAllIssuesWithWorklogsByDateRange gets ALL issues (handles pagination)
func (c *Client) GetAllIssuesWithWorklogsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]string, error) {
	logger := slog.Default()

	jql := fmt.Sprintf("worklogDate >= '%s' AND worklogDate <= '%s' AND assignee = currentUser() ORDER BY updated DESC",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	var allKeys []string
	var nextPageToken *string
	maxResults := 100
	pageNum := 1

	for {
		// Correct payload structure for /rest/api/3/search/jql
		payload := map[string]interface{}{
			"jql":        jql,
			"maxResults": maxResults,
			"fields":     []string{"key"},
		}

		// Add nextPageToken if we have one (for pagination)
		if nextPageToken != nil {
			payload["nextPageToken"] = *nextPageToken
		}

		log.Println(payload)

		var result struct {
			Issues []struct {
				Key string `json:"key"`
			} `json:"issues"`
			Total         int     `json:"total"`
			MaxResults    int     `json:"maxResults"`
			NextPageToken *string `json:"nextPageToken,omitempty"`
		}

		resp, err := c.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetBody(payload).
			SetResult(&result).
			Post("/rest/api/3/search/jql")

		if err != nil {
			logger.Error("Error calling Jira API",
				"error", err,
				"page", pageNum)
			return nil, fmt.Errorf("searching issues: %w", err)
		}

		if resp.IsError() {
			logger.Error("Jira API returned error",
				"status_code", resp.StatusCode(),
				"status", resp.Status(),
				"body", string(resp.Body()),
				"page", pageNum)

			// Check if it's authentication issue
			if resp.StatusCode() == 401 {
				return nil, fmt.Errorf("jira authentication failed - check your email and API token")
			}

			return nil, fmt.Errorf("jira API error: %s - %s", resp.Status(), string(resp.Body()))
		}

		// Add keys from this page
		for _, issue := range result.Issues {
			allKeys = append(allKeys, issue.Key)
		}

		// Check if we have more pages
		if result.NextPageToken == nil || *result.NextPageToken == "" {
			break
		}

		// Move to next page using token
		nextPageToken = result.NextPageToken
		pageNum++
	}

	return allKeys, nil
}

// Add min helper function at the bottom of the file
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LogWork creates a new worklog entry in Jira
func (c *Client) LogWork(ctx context.Context, issueKey string, timeSpentSeconds int, started time.Time, comment string) (string, error) {
	payload := map[string]interface{}{
		"timeSpentSeconds": timeSpentSeconds,
		"started":          started.Format("2006-01-02T15:04:05.000-0700"),
		"comment": map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": comment,
						},
					},
				},
			},
		},
	}

	var result Worklog
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&result).
		Post(fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueKey))

	if err != nil {
		return "", fmt.Errorf("logging work: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("jira API error: %s - %s", resp.Status(), string(resp.Body()))
	}

	return result.ID, nil
}
