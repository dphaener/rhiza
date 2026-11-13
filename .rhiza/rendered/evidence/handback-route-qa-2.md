# handback-route-qa-2

evidence · `01M7BFRKH5E07Q7WBE0VEBREK8` · actor [claude-code-420fd0ba](../actor/claude-code-420fd0ba.md)

## result

PASS

## statement

QA cycle 2, scratch build, fake rex/claude: included workflow hbi (agent build -> gate check, on.fail = handback:build) under top-level hbo: two gate reds re-entered build (same session ID) and the third pass ended the run done. With handback_cap = 1 on the top-level file the second red parked the run stopped; rhiza inbox showed a handback item naming step hbi/check and the cause with the kept gate log path, and run resume finished green. TestHandbackTarget, TestHandbackRoute, TestInboxHandbackParked pass. Rough edge: handback_cap set only in the included file is accepted by workflow check and ignored.

## Links

### Out

- verifies [handback-route](../capability/rhiza/workflows/workflow-runs/handback-route.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
