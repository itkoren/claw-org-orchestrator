#!/usr/bin/env bash
set -euo pipefail

# Uses 1Password Service Account (OP_SERVICE_ACCOUNT_TOKEN must be set in launchd)
# All secrets live in vault "Claw"

export SUPABASE_URL="$(op read "op://Claw/supabase/url")"
export SUPABASE_SERVICE_ROLE_KEY="$(op read "op://Claw/supabase/service_role_key")"

export SLACK_BOT_TOKEN="$(op read "op://Claw/slack/bot_token")"
export SLACK_SIGNING_SECRET="$(op read "op://Claw/slack/signing_secret")"

export GITHUB_TOKEN="$(op read "op://Claw/github/bot_token")"
export BRAIN_REPO_SSH="$(op read "op://Claw/github/brain_repo_ssh")"

exec /usr/local/bin/orchestrator
