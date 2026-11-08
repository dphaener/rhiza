# health-autonomy-risk-verified

evidence · `01M6W9A3FF82HQ1N1H3WAX1PJM` · actor [claude-code-bb5c79c3](../actor/claude-code-bb5c79c3.md)

## result

PASS

## statement

Verifier at 9d6f60e8: go test -trimpath ./internal/api -run TestHealthAutonomyRisk -v -count=1 passes its three subtests: flagged plan (state warn, summary has 'outside text reaches agents with no permission checks', exactly one 'workflow plan runs from a queued request with bypassed permissions in interview' detail line, Next 'rhiza guide workflows'), nothing flagged (ok, no Next, no phrase, no workflow detail), flagged while paused (warn, 'daemon paused; ' prefix plus phrase, Next stays 'rhiza daemon set rhiza.daemon.paused false'). go test -trimpath ./cmd/rhiza -run TestScripts/health -count=1 passes, including health-risk-warning: rhiza health exits 0 with 'warn autonomy ... outside text reaches agents...' and the detail line, --json shows part autonomy state warn, and acceptEdits clears it. Read the diff: healthAutonomy calls requestExposedWorkflows, adds one ExposureDetail per exposed workflow (name, trigger, bypassing steps), never sets bad, and the paused branch runs last so it keeps its Next. Full internal/api package: every test PASS; the run reported only the root guard's .rhiza-root/usage.live.json change, written by the live root daemon, not by this diff.

## Links

### Out

- verifies [health-autonomy-risk](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-autonomy-risk.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
