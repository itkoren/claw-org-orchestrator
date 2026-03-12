package scheduler

import (
	"fmt"

	"github.com/itkoren/claw-org-orchestrator/internal/config"
	"github.com/itkoren/claw-org-orchestrator/internal/slack"
	"github.com/itkoren/claw-org-orchestrator/internal/supabase"
)

func Tick(cfg *config.Config, sb *supabase.Client, sl *slack.Client) error {
	// 1) Pull due schedules -> instantiate tasks
	// 2) Pull queued tasks -> delegate/execute
	// 3) Push status updates to Slack threads
	// 4) Persist run events to Supabase
	// 5) Trigger brain sync if needed

	due, err := sb.ListDueSchedules()
	if err != nil {
		return err
	}
	for _, s := range due {
		taskID, err := sb.InstantiateTaskFromSchedule(s)
		if err != nil {
			_ = sl.PostOrgAnnouncement(fmt.Sprintf("Schedule instantiate failed: %v", err))
			continue
		}
		_ = sl.PostOrgAnnouncement(fmt.Sprintf("Scheduled task created: %s", taskID))
	}

	queued, err := sb.ListQueuedTasks()
	if err != nil {
		return err
	}
	for _, t := range queued {
		// Delegation rule:
		// immediate=true => Satya executes
		// else Satya delegates to team managers (Elon for engineering)
		if err := RunTask(cfg, sb, sl, t.ID); err != nil {
			_ = sb.MarkTaskFailed(t.ID, err.Error())
		}
	}
	return nil
}
