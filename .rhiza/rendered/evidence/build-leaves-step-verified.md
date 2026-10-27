# build-leaves-step-verified

evidence · `01M5ZMT4JXX7W9SHT0BA17T8SN` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Parse clauses: TestLeavesStepLoads loads a leaves step with both builder and verifier sub-tables (split, interval, per-role models); TestLeavesStepRefused rejects a missing builder, a missing verifier, an unknown sub-table key (retries), models on wait, command and leaves steps, split on an agent step, and split = "diagonal". Engine clause: TestLeavesStepTwoLeaves runs a fake runner over two leaves and finds a leaf-session, leaf-attempt and leaf outcome entry per leaf and role in the journal. Status clause: testscript build-leaves-status finds rhiza run status listing each leaf with its builder and verifier sessions, judged attempts and outcome, in text and --json. All passed.

## tests

go test -count=1 ./internal/workflow -run 'TestLeavesStepLoads|TestLeavesStepRefused'; ./internal/run -run TestLeavesStepTwoLeaves; ./cmd/rhiza -run 'TestScripts/build-leaves-status$' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-leaves-step](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
