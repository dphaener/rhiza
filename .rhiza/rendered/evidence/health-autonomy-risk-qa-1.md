# health-autonomy-risk-qa-1

evidence · `01M6W9QDKBP07NEKV9TBJXEW7R` · actor [claude-code-afe51521](../actor/claude-code-afe51521.md)

## result

PASS

## statement

QA at bbfe87e3 with a built binary (go build -trimpath) in a scratch Rhiza root and project, real git config, no fakes: rhiza health with defaults: autonomy ok, no phrase. After git config rhiza.autonomy.plan auto and rhiza.daemon.permissionMode bypassPermissions: exit 0, 'warn autonomy plan auto, accept auto, ship auto, signoff auto; outside text reaches agents with no permission checks', detail lines 'workflow plan runs from a queued request with bypassed permissions in interview, review, critic, revise, arbiter' and the ship-change one (accept/ship default auto), next rhiza guide workflows; --json autonomy state warn. With rhiza.daemon.paused true: 'daemon paused; ...' plus the phrase, next stays rhiza daemon set rhiza.daemon.paused false. permissionMode acceptEdits: ok, no phrase, no workflow lines.

## Links

### Out

- verifies [health-autonomy-risk](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-autonomy-risk.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
