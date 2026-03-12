package supabase

type Task struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	Immediate     bool   `json:"immediate"`
	SlackChannelID string `json:"slack_channel_id"`
	SlackThreadTS  string `json:"slack_thread_ts"`
	OwnerTeamID    string `json:"owner_team_id"`
	OwnerAgentID   string `json:"owner_agent_id"`
	CreatedAt      string `json:"created_at"`
}
