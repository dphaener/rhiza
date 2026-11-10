# custom-workflows-doc-voice-qa-1

evidence · `01M73F1ETMTC9RGHHA7H0CKYBQ` · actor [claude-code-4bad3217](../actor/claude-code-4bad3217.md)

## result

PASS

## statement

QA at 89bc7981 with a build of the branch in a scratch repo: rhiza workflow new bugfix and rhiza workflow check bugfix succeed; rhiza change new "Totals are wrong on the invoice page" --workflow bugfix exits 0 and rhiza change show lists Workflow: bugfix; rhiza change workflow --clear prints that the dispatch rules choose again, and rhiza change workflow bugfix pins it again, all exit 0. go test -trimpath ./cmd/rhiza -run 'TestScripts/(custom-workflows-doc|docs-)' passes (GOTMPDIR outside the repo). The owner's approval is left to sign-off.

## Links

### Out

- verifies [custom-workflows-doc-voice](../capability/rhiza/documentation/custom-workflows-doc-voice.md)

## History

- `01M71G4MMC4VGT41WVM4A33ZWM` Agents, runs and workflows pages in the owner's voice: docs/agents.md, runs-and-daemon.md and custom-workflows.md rewritten with humanize, examples as plain commands, signed off by the owner
