# runs-daemon-doc-voice-qa-1

evidence · `01M73F1EJTQZSEN83EACHANM60` · actor [claude-code-4bad3217](../actor/claude-code-4bad3217.md)

## result

PASS

## statement

QA at 89bc7981 with a build of the branch: rhiza run send-back --help shows rhiza run send-back [<run>] --note <note>, matching the page's new example; in a scratch repo rhiza run send-back --note "The retry loop needs a cap" and rhiza run approve both parse and exit 2 with 'works on a run waiting at a sign-off step: no runs yet', as documented. rhiza run request-changes --help and rhiza describe run.request-changes --json exist as the page's fallback paragraph says. go test -trimpath ./internal/doclinks/... ok and make self-check in a clean worktree at HEAD exits 0. The owner's approval is left to sign-off.

## Links

### Out

- verifies [runs-daemon-doc-voice](../capability/rhiza/documentation/runs-daemon-doc-voice.md)

## History

- `01M71G4MMC4VGT41WVM4A33ZWM` Agents, runs and workflows pages in the owner's voice: docs/agents.md, runs-and-daemon.md and custom-workflows.md rewritten with humanize, examples as plain commands, signed off by the owner
