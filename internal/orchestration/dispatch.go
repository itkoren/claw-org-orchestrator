package orchestration

import (
	"fmt"

	"github.com/itkoren/claw-org-orchestrator/internal/brain"
	"github.com/itkoren/claw-org-orchestrator/internal/config"
	"github.com/itkoren/claw-org-orchestrator/internal/routing"
	"github.com/itkoren/claw-org-orchestrator/internal/slack"
	"github.com/itkoren/claw-org-orchestrator/internal/supabase"
)

func ProcessTask(cfg *config.Config, sb *supabase.Client, sl *slack.Client, taskID string) error {
	t, err := sb.GetTask(taskID)
	if err != nil {
		return err
	}

	// Mark running early to avoid double-processing in v1 polling.
	if err := sb.MarkTaskRunning(t.ID); err != nil {
		return err
	}

	// Ensure Slack thread exists.
	channel := t.SlackChannelID
	if channel == "" {
		// v1 default: post everything to org channel if no owner channel is set yet
		channel = cfg.SlackOrgChannelID
	}
	threadTS := t.SlackThreadTS
	if threadTS == "" {
		rootTS, err := sl.PostRoot(channel, fmt.Sprintf("*TASK* %s\n%s", t.Title, t.Description))
		if err != nil {
			return err
		}
		threadTS = rootTS
		_ = sb.SetSlackThread(t.ID, channel, threadTS)
	}

	decision := routing.DecideExecutor(t.Immediate, t.Title, t.Description)
	_ = sl.Reply(channel, threadTS, fmt.Sprintf("Satya routing: `%s` (%s)", decision.ExecutorAgent, decision.Reason))

	// Call runner (OpenClaw sidecar) with the chosen executor as “agentName”.
	runner := NewRunnerClient(cfg.RunnerURL)
	resp, err := runner.Run(RunRequest{
		AgentName: decision.ExecutorAgent,
		ModelTier: "tier2",
		TaskID:    t.ID,
		Title:     t.Title,
		Prompt:    t.Description,
	})
	if err != nil {
		_ = sl.Reply(channel, threadTS, "Execution error: "+err.Error())
		return err
	}

	_ = sl.Reply(channel, threadTS, "*Output*\n"+resp.Output)

	// Optional: write artifacts to brain + sync
	if cfg.BrainRepoDir != "" {
		_ = brain.WriteTaskArtifact(cfg.BrainRepoDir, t.ID, "output.md", resp.Output)
		_ = brain.Sync(cfg.BrainRepoDir)
	}

	if err := sb.MarkTaskDone(t.ID); err != nil {
		return err
	}
	_ = sl.Reply(channel, threadTS, "Status: `done`")
	return nil
}
