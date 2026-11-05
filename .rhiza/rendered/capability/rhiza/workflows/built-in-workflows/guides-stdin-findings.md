# guides-stdin-findings

capability · `01M6MMAF8Z4XWD1JD9WG5PRYBR` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

Tests in internal/guide (build_test.go, review_test.go, plan_test.go) assert each of the three guides shows the `- <<'EOF'` form for every rhiza run findings and rhiza run verdict it names and that no `<file>` form of either command remains in them; go test -trimpath ./internal/guide/... passes

## statement

The build, review and plan-review guides show rhiza run findings and rhiza run verdict fed from stdin by a heredoc (`- <<'EOF'` ... `EOF`), with their existing flags (--leaf, --check, --json) and JSON shapes, so an agent records findings and verdicts without a temp file or compound shell command

## Links

### Out

- decided-by [stdin-guides-ruling](../../../../decision/stdin-guides-ruling.md)

### In

- verifies [guides-stdin-findings-verified](../../../../evidence/guides-stdin-findings-verified.md)
- verifies [guides-stdin-findings-qa-1](../../../../evidence/guides-stdin-findings-qa-1.md)
- verifies [guides-stdin-findings-verified-2](../../../../evidence/guides-stdin-findings-verified-2.md)
- verifies [pr-282-fix-agents-in-runs-hit-claude](../../../../evidence/pr-282-fix-agents-in-runs-hit-claude.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
