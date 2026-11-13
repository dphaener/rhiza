# run-handback-qa-2

evidence · `01M7BFRKN84WG1KQZY220Q0QVH` · actor [claude-code-420fd0ba](../actor/claude-code-420fd0ba.md)

## result

PASS

## statement

QA cycle 2, scratch build, fake claude: a red gate step handed back twice; each resume (claude --resume, same session ID) got the handback text 'Fix it, see <session dir>/handback.md', and handback.md named failing step hbi/check, outcome fail and '`gatecheck`: output in <run dir>/gate-hbi_check-<visit>-2.log' for the matching visit. TestHandbackDeliver and TestScripts/run-handback pass.

## Links

### Out

- verifies [run-handback](../capability/rhiza/workflows/workflow-runs/run-handback.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
