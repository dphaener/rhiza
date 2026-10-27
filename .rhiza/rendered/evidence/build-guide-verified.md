# build-guide-verified

evidence · `01M5ZMT81G820REYXB6F49FHT1` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Testscript guide finds rhiza guide ls listing the build topic and rhiza guide build carrying the Leaf: trailer, the outside-work rule (blocks question, discovered-from draft) and the Evidence-or-findings rule; TestBuildText checks the same phrases. guide-install and TestGuideInstallBuild find rhiza guide install writing skills/rhiza-build/SKILL.md whose body equals rhiza guide build. All passed.

## tests

go test -count=1 ./internal/guide -run TestBuildText; ./internal/api -run TestGuideInstallBuild; ./cmd/rhiza -run 'TestScripts/(guide|guide-install)$' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-guide](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
