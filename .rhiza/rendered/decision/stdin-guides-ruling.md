# stdin-guides-ruling

decision · `01M6MM7RZ3R95169W99F1C8BZ7`

## reasoning

The change set's intent is to stop runs stalling on Claude Code permission prompts, which blocks the daemon shipping the backlog unattended; recent run-reliability fixes (#259, #267) serve milestone-25 for the same reason. The guides (internal/guide/text/build.md:78, review.md:37 and :50, plan-review.md:33 and :45) show `<file>` for both findings and verdict, and verdict is written the same way findings are, so leaving it would keep the same prompt-triggering pattern; switching both is the whole fix. readReviewFile (internal/api/run_review.go:276) already reads `-` as stdin and run-review-journal.txtar exercises it, so the only change and the only new test surface is the guide text, which internal/guide's existing tests already pin by substring.

## statement

Change set 01M6MM3 serves milestone-25 (the daemon builds and ships the backlog on its own). Every place the build, review and plan-review guides show `rhiza run findings ... <file>` or `rhiza run verdict ... <file>` switches to the stdin heredoc form (e.g. `rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'` ... `EOF`), keeping the existing flags (--leaf, --check, --json) and JSON shapes; no temp file or compound shell is shown anywhere. No CLI change: both commands already accept `-`. Test it in internal/guide: update build_test.go, review_test.go and plan_test.go to assert the `- <<'EOF'` form for findings and verdict and that no `<file>` form for either command remains in the three guides; run `go test -trimpath ./internal/guide/...`. Stdin reading itself is already covered by run-review-journal.txtar.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-9bf069b9-1](../question/question-9bf069b9-1.md)
- derived-from [stdin-findings-guides-plan](../source/stdin-findings-guides-plan.md)

### In

- decided-by [guides-stdin-findings](../capability/rhiza/workflows/built-in-workflows/guides-stdin-findings.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
