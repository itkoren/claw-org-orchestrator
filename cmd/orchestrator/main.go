package main

import (
	"log"
	"time"

	"github.com/itkoren/claw-org-orchestrator/internal/config"
	"github.com/itkoren/claw-org-orchestrator/internal/scheduler"
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
		if err := scheduler.Tick(cfg, sb, sl); err != nil {
			log.Printf("tick error: %v\n", err)
		}
	}
}
