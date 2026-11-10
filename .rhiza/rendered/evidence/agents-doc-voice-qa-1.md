# agents-doc-voice-qa-1

evidence · `01M73F1E8XQ76ZV2B2P4E53AT9` · actor [claude-code-4bad3217](../actor/claude-code-4bad3217.md)

## result

PASS

## statement

QA at 89bc7981 with a build of the branch in a scratch git repo (RHIZA_* unset, own HOME): ran the docs/agents.md commands by hand. rhiza guide ls and rhiza guide build print the topics; rhiza guide install --agent claude writes the skills to ~/.claude/skills/ and --project to .claude/skills/. The by-hand example rhiza add question export-size-cap --field statement="Do exports need a size cap?" exits 0 and writes the node into the current change set; the agent write-back path rhiza apply <file> --json with an add write exits 0 with a JSON entry. grep finds rhiza apply <file> --json in step 3 of the loop. go test -trimpath ./internal/doclinks/... ok, and make self-check in a clean worktree at HEAD exits 0 (warns only on old unclaimed main commits). The owner's approval is left to sign-off.

## Links

### Out

- verifies [agents-doc-voice](../capability/rhiza/documentation/agents-doc-voice.md)

## History

- `01M71G4MMC4VGT41WVM4A33ZWM` Agents, runs and workflows pages in the owner's voice: docs/agents.md, runs-and-daemon.md and custom-workflows.md rewritten with humanize, examples as plain commands, signed off by the owner
