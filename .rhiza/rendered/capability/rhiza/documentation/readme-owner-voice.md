# readme-owner-voice

capability · `01M71G7PH3S3AB0X96JD4H94BY` · parent [documentation](../documentation.md)

## acceptance

go test -trimpath ./internal/doclinks/... passes (no broken relative link in the real README.md and docs/) and make self-check passes (rhiza check's docs-commands rule finds no rhiza command in README.md or docs/ that does not exist), README.md no longer says pages are on the way, and the owner approves the page at the run's sign-off

## statement

README.md reads in the owner's voice (drafted with the humanize skill) and its where-next section links the docs that exist (quickstart, concepts, agents, runs and the daemon, custom workflows, the generated reference) instead of saying they are on the way

## Links

### Out

- decided-by [m31-polish-request](../../../decision/m31-polish-request.md)

### In

- verifies [readme-owner-voice-verified](../../../evidence/readme-owner-voice-verified.md)
- verifies [readme-owner-voice-qa-1](../../../evidence/readme-owner-voice-qa-1.md)
- verifies [readme-owner-voice-verified-2](../../../evidence/readme-owner-voice-verified-2.md)
- verifies [pr-427-readme-and-quickstart-polish](../../../evidence/pr-427-readme-and-quickstart-polish.md)

## History

- `01M71G4BKFHP1E172RR6H82WFT` README and quickstart in the owner's voice: rewritten with humanize, the quickstart writes the model with plain rhiza add/link commands instead of an apply JSON file, the README's stale 'on the way' lines go, signed off by the owner
