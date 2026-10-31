# command-step-keeps-output-verified

evidence · `01M69NWH1VF4CN7T9GXQCBF6R0` · actor [claude-code-f7a86a7b](../actor/claude-code-f7a86a7b.md)

## result

PASS

## statement

At 17514df9: go test ./cmd/rhiza -run TestScripts/run-command-step-output passes; the script runs 'echo out1; echo boom >&2; exit 3' and run status --json shows detail "exit 3\nout1\nboom", a 500-line step keeps lines 481-500 (line 500 present, line 1 and line 480 absent), and a passing step's detail is exactly "exit 0". go test ./internal/run -run TestCommandDetailKeepsRunnerOutput passes, showing a runner's detail kept under the exit line. Diff read: stepRunner.Command merges stderr into stdout (exec 2>&1) and keeps a tail capped at 20 lines and 2048 bytes, cut on a UTF-8 rune start; detail is returned only for non-zero exits. go test ./... passes except a timing test (modelapp TestGeneratedModel, 55ms vs 50ms) that passes when rerun alone; make self-check passes with RHIZA_ROOT unset.

## Links

### Out

- verifies [command-step-keeps-output](../capability/rhiza/workflows/workflow-definitions/command-step/command-step-keeps-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
