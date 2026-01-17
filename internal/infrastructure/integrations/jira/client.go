// internal/infrastructure/integrations/jira/client.go
package jira

import (
	"context"
	"fmt"
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
	ID               string    `json:"id"`
	IssueKey         string    `json:"issueKey"`
	TimeSpent        string    `json:"timeSpent"`
	TimeSpentSeconds int       `json:"timeSpentSeconds"`
	Started          time.Time `json:"started"`
	Comment          string    `json:"comment"`
	Author           Author    `json:"author"`
}

type Author struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

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
		return nil, fmt.Errorf("jira API error: %s", resp.Status())
	}

	return result.Worklogs, nil
}

func (c *Client) LogWork(ctx context.Context, issueKey string, timeSpentSeconds int, started time.Time, comment string) (string, error) {
	payload := map[string]interface{}{
		"timeSpentSeconds": timeSpentSeconds,
		"started":          started.Format("2006-01-02T15:04:05.000-0700"),
		"comment":          comment,
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
		return "", fmt.Errorf("jira API error: %s", resp.Status())
	}

	return result.ID, nil
}
