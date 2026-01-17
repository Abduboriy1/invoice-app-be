package jobs

import (
	"context"
	"log"

	"github.com/invoice-app-be/internal/infrastructure/integrations/jira"
)

type SyncJiraJob struct {
	syncer *jira.Syncer
}

func NewSyncJiraJob(syncer *jira.Syncer) *SyncJiraJob {
	return &SyncJiraJob{syncer: syncer}
}

func (j *SyncJiraJob) Run(ctx context.Context) error {
	log.Println("Starting Jira sync job...")

	// TODO: Implement job logic
	// - Fetch users with Jira integration enabled
	// - For each user, sync their worklogs

	log.Println("Jira sync job completed")
	return nil
}
