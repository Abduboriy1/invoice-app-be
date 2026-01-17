// internal/infrastructure/database/postgres/timeentry_repository.go
package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/invoice-app-be/internal/domain/timeentry"
)

type TimeEntryRepository struct {
	db *sqlx.DB
}

func NewTimeEntryRepository(db *sqlx.DB) *TimeEntryRepository {
	return &TimeEntryRepository{db: db}
}

func (r *TimeEntryRepository) Create(ctx context.Context, entry *timeentry.TimeEntry) error {
	query := `
        INSERT INTO time_entries (id, user_id, invoice_id, description, hours, hourly_rate, date,
                                jira_issue_key, jira_worklog_id, jira_synced_at, is_billable, is_invoiced,
                                created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `
	_, err := r.db.ExecContext(ctx, query, entry.ID, entry.UserID, entry.InvoiceID, entry.Description,
		entry.Hours, entry.HourlyRate, entry.Date, entry.JiraIssueKey, entry.JiraWorklogID,
		entry.JiraSyncedAt, entry.IsBillable, entry.IsInvoiced, entry.CreatedAt, entry.UpdatedAt)
	return err
}

func (r *TimeEntryRepository) GetByID(ctx context.Context, id uuid.UUID) (*timeentry.TimeEntry, error) {
	var entry timeentry.TimeEntry
	query := `SELECT * FROM time_entries WHERE id = $1`
	if err := r.db.GetContext(ctx, &entry, query, id); err != nil {
		return nil, fmt.Errorf("getting time entry: %w", err)
	}
	return &entry, nil
}

func (r *TimeEntryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]timeentry.TimeEntry, error) {
	var entries []timeentry.TimeEntry
	query := `SELECT * FROM time_entries WHERE user_id = $1 ORDER BY date DESC`
	if err := r.db.SelectContext(ctx, &entries, query, userID); err != nil {
		return nil, fmt.Errorf("getting time entries: %w", err)
	}
	return entries, nil
}

func (r *TimeEntryRepository) GetByJiraWorklogID(ctx context.Context, worklogID string) (*timeentry.TimeEntry, error) {
	var entry timeentry.TimeEntry
	query := `SELECT * FROM time_entries WHERE jira_worklog_id = $1`
	if err := r.db.GetContext(ctx, &entry, query, worklogID); err != nil {
		return nil, fmt.Errorf("getting time entry by jira worklog: %w", err)
	}
	return &entry, nil
}

func (r *TimeEntryRepository) Update(ctx context.Context, entry *timeentry.TimeEntry) error {
	query := `
        UPDATE time_entries SET description = $2, hours = $3, jira_worklog_id = $4, updated_at = $5
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, entry.ID, entry.Description, entry.Hours, entry.JiraWorklogID, entry.UpdatedAt)
	return err
}

func (r *TimeEntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM time_entries WHERE id = $1", id)
	return err
}
