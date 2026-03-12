package main

import (
	"log"
	"time"

	"github.com/itkoren/claw-org-orchestrator/internal/config"
	"github.com/itkoren/claw-org-orchestrator/internal/orchestration"
	"github.com/itkoren/claw-org-orchestrator/internal/slack"
	"github.com/itkoren/claw-org-orchestrator/internal/supabase"
)

func main() {
	cfg := config.MustLoad()

	sb := supabase.MustNew(cfg)
	sl := slack.MustNew(cfg)

	log.Println("openclaw orchestrator starting...")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Boot message (optional)
	_ = sl.PostOrgAnnouncement("Orchestrator online (Mac mini).")

	for {
		<-ticker.C
		queued, err := sb.ListQueuedTasks(25)
		if err != nil {
			log.Printf("ListQueuedTasks error: %v\n", err)
			continue
		}
		for _, t := range queued {
			if err := orchestration.ProcessTask(cfg, sb, sl, t.ID); err != nil {
				log.Printf("ProcessTask(%s) error: %v\n", t.ID, err)
				_ = sb.MarkTaskFailed(t.ID, err.Error())
			}
		}
	}
}
