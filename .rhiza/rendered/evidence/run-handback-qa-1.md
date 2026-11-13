# run-handback-qa-1

evidence · `01M7BEHN0YVWS4HMCXAFN5HVEB` · actor [claude-code-a9fe07bb](../actor/claude-code-a9fe07bb.md)

## result

PASS

## statement

QA cycle 1, scratch build: the resumed agent session read handback.md: for a red command step it named `check` and the output tail 'tests red 3'; for a red gate step it named `check` with 'output in <run dir>/gate-gate-1-2.log'. Each resume carried the step's handback text ('Fix it, see <session dir>/handback.md') with claude --resume and the same session ID. A session whose transcript is gone failed the step with 'session <id> cannot be resumed: no transcript ...'; an agent step with no handback template failed it naming the step and session. Scoped tests TestHandbackDeliver and TestScripts/run-handback pass.

## Links

### Out

- verifies [run-handback](../capability/rhiza/workflows/workflow-runs/run-handback.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
