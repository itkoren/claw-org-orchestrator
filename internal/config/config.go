package config

import (
	"log"
	"os"
)

type Config struct {
	SupabaseURL            string
	SupabaseServiceRoleKey string

	SlackBotToken     string
	SlackOrgChannelID string // e.g. C0123... for #org-command-center

	RunnerURL string // http://127.0.0.1:8787

	BrainRepoDir string // local path to checked out openclaw-brain (optional v1)
}

func MustLoad() *Config {
	cfg := &Config{
		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseServiceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		SlackBotToken:          os.Getenv("SLACK_BOT_TOKEN"),
		SlackOrgChannelID:      os.Getenv("SLACK_ORG_CHANNEL_ID"),
		RunnerURL:              getenv("RUNNER_URL", "http://127.0.0.1:8787"),
		BrainRepoDir:           getenv("BRAIN_REPO_DIR", ""),
	}

	must(cfg.SupabaseURL, "SUPABASE_URL")
	must(cfg.SupabaseServiceRoleKey, "SUPABASE_SERVICE_ROLE_KEY")
	must(cfg.SlackBotToken, "SLACK_BOT_TOKEN")
	must(cfg.SlackOrgChannelID, "SLACK_ORG_CHANNEL_ID")
	return cfg
}

func must(v, name string) {
	if v == "" {
		log.Fatalf("missing env: %s", name)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
