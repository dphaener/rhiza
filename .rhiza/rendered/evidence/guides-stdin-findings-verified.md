# guides-stdin-findings-verified

evidence · `01M6MT9MD2RQDPES6QGDTBW357` · actor [claude-code-95a2eb87](../actor/claude-code-95a2eb87.md)

## result

PASS

## statement

At c0ce8402 ran go test -trimpath -count=1 ./internal/guide/... (ok). build_test.go, review_test.go and plan_test.go assert the `- <<'EOF'` heredoc form (and closing EOF) for each rhiza run findings / rhiza run verdict the build, review and plan-review guides name, and TestGuidesFeedFindingsFromStdin fails on any remaining `<file>` form. grep of the three guide texts shows five commands, all heredoc on stdin with their existing flags (--leaf, --check, --json) and the unchanged JSON shapes. readReviewFile in internal/api/run_review.go reads stdin for `-` for both findings and verdict, so the shown commands work.

## Links

### Out

- verifies [guides-stdin-findings](../capability/rhiza/workflows/built-in-workflows/guides-stdin-findings.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
