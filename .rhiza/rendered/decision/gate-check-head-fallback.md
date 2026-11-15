# gate-check-head-fallback

decision · `01M7EHQ2GQYWMFNQZH58BD3X26`

## reasoning

The orchestrator's hint asks for a fallback to the installed rhiza only when no head build exists. Land and render already treat a named-but-unbuildable head as an error (conductor-land-render-head-build), and checking with the installed binary in that case would bring back exactly the false docs/reference failures this change set fixes, or worse a false pass. shipRunHead returns an error on any non-zero exit, but check exits 1 on findings, so the report must be parsed from the child's stdout as prePushHeadCheck does.

## statement

rhiza check falls back to the running (installed) rhiza only when no head build is named: SHIP_HEAD_RHIZA unset, naming the running binary, or a running binary that is not the rhiza PATH resolves (make self-check's fresh build, which must check with this commit's build). A named head that is missing is built on demand under SHIP_HEAD_BUILD (shipHeadBuild); a missing head with SHIP_HEAD_BUILD empty, a failed build, or a head whose output is not a check report fails the check with an error naming the head build, never a silent in-process check. A head report with failing findings exits non-zero, so check reads its --json report even when the child exits non-zero and reports those findings with check's usual exit.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [gate-check-head-build-plan](../source/gate-check-head-build-plan.md)

### In

- decided-by [gate-check-head-build](../capability/rhiza/orchestration/shipping/gate-check-head-build.md)

## History

- `01M7EHHAG3CTRQCJCWSG2W7BTZ` Fix the gate step: the conductor default's gate runs rhiza check with the installed rhiza, not the PR head's build (SHIP_HEAD_RHIZA), so a change set that adds a CLI command, flag or setting fails the gate on docs/reference files its own build passes and hands back to conduct in a loop (runs 01M7DVV7GNJAWG29QZC5F2JWC3 and 01M7E56MGTY4M60XZNEG2A6537; request 01M7EAYSYK7Z32QRP0FB4XPKJN)
