#!/usr/bin/env bash
set -euo pipefail

# 1) Dependencies
brew install go node pnpm git 1password-cli supabase/tap/supabase vercel-cli duplicati

# 2) Build Go orchestrator
mkdir -p /Users/claw/bin
go build -o /Users/claw/bin/orchestrator ./cmd/orchestrator

# 3) Install runner
cd runner/claw-runner
pnpm install
pnpm build

echo "Next: configure launchd + start runner (use pm2 or launchd)."
