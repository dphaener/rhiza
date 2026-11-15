# gate-check-via-head

decision · `01M7EHQ2F94B5FTJF7JGPVWXCJ`

## reasoning

The gate step runs its commands with sh -c in the run worktree (stepRunner.CommandOutput), so rhiza there is the installed binary, while the step's environment already carries the workflow's [env] (stepEnv), SHIP_HEAD_RHIZA included. conductor-land-render-head-build solved the same problem for land and render inside the operations with runViaHead, and ship open-pr's pre-push check builds the head for its check (prePushHeadCheck). Changing the builtin's gate command to a shell fallback such as ${SHIP_HEAD_RHIZA:-rhiza} would skip the on-demand build (shipHeadBuild) and put Rhiza-repo knowledge into the generic builtin. Redirecting inside runCheck itself would also redirect the pre-push check's in-process call in its temporary checkout, which already has its own head path, so the redirect goes at the operation's entry.

## statement

rhiza check itself runs through the run's head build, as rhiza land and rhiza render already do (runViaHead): when the environment names SHIP_HEAD_RHIZA, it is not the running binary, and the running binary is the rhiza PATH resolves (as in the gate's plain rhiza check; a fresh build run by path, such as make self-check's, checks in-process), the check operation execs that build's check --json in the same directory with the request's flags (--base, --skip-reference, --branch-newer-format) and reports its result. The builtin ship-change gate's commands stay as they are (plain rhiza check), so every workflow whose [env] names a head build gets it with no workflow change. Only the check operation (CLI and api call) does this; the in-process callers of runCheck (the pre-push check, run code-check, health) keep their own behaviour.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [gate-check-head-build-plan](../source/gate-check-head-build-plan.md)

### In

- decided-by [gate-check-head-build](../capability/rhiza/orchestration/shipping/gate-check-head-build.md)

## History

- `01M7EHHAG3CTRQCJCWSG2W7BTZ` Fix the gate step: the conductor default's gate runs rhiza check with the installed rhiza, not the PR head's build (SHIP_HEAD_RHIZA), so a change set that adds a CLI command, flag or setting fails the gate on docs/reference files its own build passes and hands back to conduct in a loop (runs 01M7DVV7GNJAWG29QZC5F2JWC3 and 01M7E56MGTY4M60XZNEG2A6537; request 01M7EAYSYK7Z32QRP0FB4XPKJN)
