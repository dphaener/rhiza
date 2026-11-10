# runs-daemon-doc-voice

capability · `01M71G7QMC5NKK7T0TNKBPKFQ7` · parent [documentation](../documentation.md)

## acceptance

go test -trimpath ./internal/doclinks/... passes (no broken relative link in the real README.md and docs/) and make self-check passes (rhiza check's docs-commands rule finds no rhiza command in README.md or docs/ that does not exist), and the owner approves docs/runs-and-daemon.md at the run's sign-off

## statement

docs/runs-and-daemon.md reads in the owner's voice, and its examples prefer plain commands (rhiza run send-back --note over a findings JSON file where either works)

## Links

### Out

- decided-by [m31-polish-request](../../../decision/m31-polish-request.md)

### In

- verifies [runs-daemon-doc-voice-verified](../../../evidence/runs-daemon-doc-voice-verified.md)
- verifies [runs-daemon-doc-voice-qa-1](../../../evidence/runs-daemon-doc-voice-qa-1.md)
- verifies [pr-428-agents-and-runs-docs-polish](../../../evidence/pr-428-agents-and-runs-docs-polish.md)

## History

- `01M71G4MMC4VGT41WVM4A33ZWM` Agents, runs and workflows pages in the owner's voice: docs/agents.md, runs-and-daemon.md and custom-workflows.md rewritten with humanize, examples as plain commands, signed off by the owner
