# custom-workflows-doc-voice-verified

evidence · `01M73DRTASYKBX7BBRPZ8WAJTB` · actor [claude-code-d4997ce5](../actor/claude-code-d4997ce5.md)

## result

PASS

## statement

At e539a7b4 ran go test -trimpath ./cmd/rhiza -run 'TestScripts/custom-workflows-doc|TestScripts/docs-' -count=1 (GOTMPDIR moved outside the repo; the run's GOTMPDIR under .git puts $WORK inside the real repo and makes docs-commands-exist's rhiza init fail regardless of the change): custom-workflows-doc and all docs-* scripts PASS. Read the diff to docs/custom-workflows.md: a few lines loosened in plain voice, plus a new 'Pin one change set to it' section naming rhiza change new --workflow bugfix, rhiza change workflow bugfix and rhiza change workflow --clear. Those text blocks aren't run by the doc test, so I ran them in a scratch repo with a build of this commit: change new --workflow, change workflow <name> and --clear all exit 0 with matching output, and the help text confirms a pin overrides the dispatch rules as the page says. The commit also stages only its own actor record. The owner's approval of the page is left to the run's sign-off.

## Links

### Out

- verifies [custom-workflows-doc-voice](../capability/rhiza/documentation/custom-workflows-doc-voice.md)

## History

- `01M71G4MMC4VGT41WVM4A33ZWM` Agents, runs and workflows pages in the owner's voice: docs/agents.md, runs-and-daemon.md and custom-workflows.md rewritten with humanize, examples as plain commands, signed off by the owner
