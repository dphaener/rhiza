# failed-step-detail-shown-qa-1

evidence · `01M69QFS6AY9TXXDC8RCGYFC0N` · actor [claude-code-58a35149](../actor/claude-code-58a35149.md)

## result

PASS

## statement

With the built binary in a scratch repo, plain rhiza run status after the out workflow failed printed '- b (visit 1): fail, then end: exit 3' then '    out1' and '    boom' indented beneath it; the run start output ends with the same block. The Rex OSC 7501 message rhiza emitted decoded to the single line 'run failed: step "b" ended with outcome "fail"', and run.go's Rex report still cuts detail at its first newline.

## Links

### Out

- verifies [failed-step-detail-shown](../capability/rhiza/workflows/workflow-runs/run-status/failed-step-detail-shown.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
