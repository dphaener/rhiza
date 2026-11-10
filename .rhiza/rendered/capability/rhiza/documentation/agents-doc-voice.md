# agents-doc-voice

capability · `01M71G7QJ4GW14BQ3D485XGYZ1` · parent [documentation](../documentation.md)

## acceptance

go test -trimpath ./internal/doclinks/... passes (no broken relative link in the real README.md and docs/) and make self-check passes (rhiza check's docs-commands rule finds no rhiza command in README.md or docs/ that does not exist), grep finds 'rhiza apply' in the agent loop's write-back step of docs/agents.md, and the owner approves docs/agents.md at the run's sign-off

## statement

docs/agents.md reads in the owner's voice; its agent loop keeps writing back with one rhiza apply <file> --json, the batch and agent path that constraint doc-examples-plain-commands and the rhiza skill prescribe, and only an example of a person writing the model by hand uses plain rhiza add, link or set commands

## Links

### Out

- decided-by [m31-polish-request](../../../decision/m31-polish-request.md)

### In

- verifies [agents-doc-voice-verified](../../../evidence/agents-doc-voice-verified.md)
- verifies [agents-doc-voice-qa-1](../../../evidence/agents-doc-voice-qa-1.md)
- verifies [pr-428-agents-and-runs-docs-polish](../../../evidence/pr-428-agents-and-runs-docs-polish.md)

## History

- `01M71G4MMC4VGT41WVM4A33ZWM` Agents, runs and workflows pages in the owner's voice: docs/agents.md, runs-and-daemon.md and custom-workflows.md rewritten with humanize, examples as plain commands, signed off by the owner
