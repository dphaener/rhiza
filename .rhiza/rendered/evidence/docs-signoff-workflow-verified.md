# docs-signoff-workflow-verified

evidence · `01M7325VNM06SGG8YPH7N638WT` · actor [claude-code-3efb2c24](../actor/claude-code-3efb2c24.md)

## result

PASS

## statement

At fbda8fda, built cmd/rhiza (-trimpath) and ran it in this checkout with RHIZA_ROOT unset: `rhiza workflow check docs` reports workflow docs is valid; `rhiza workflow describe docs` prints 'autonomy signoff owner' and --json gives autonomy [{"mode":"owner","stage":"signoff"}] (fixing attempt 1's finding); `change show --json` of 01M71G4BKFHP1E172RR6H82WFT, 01M71G4G92NEJ33C04DSKC4RMK and 01M71G4MMC4VGT41WVM4A33ZWM (README/quickstart, Concepts, Agents and runs pages) each shows workflow docs. docs.toml matches builtin ship-change apart from name, [autonomy] signoff = owner and the docs-build/docs-review includes, whose builder and fixer prompts load the humanize skill under constraints doc-prose-humanized and doc-examples-plain-commands. go test ./cmd/rhiza -run TestScripts/(workflow-describe|describe-golden) and go test ./internal/api -run Workflow pass. make self-check passes in a clean worktree at the commit with the run's untracked entry 01M72YRNA8733AQFN5R49WHEXE added (in the run checkout it fails only on untracked run/actor entries, not on the commit).

## Links

### Out

- verifies [docs-signoff-workflow](../capability/rhiza/documentation/docs-signoff-workflow.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
