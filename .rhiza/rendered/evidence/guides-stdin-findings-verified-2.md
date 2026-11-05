# guides-stdin-findings-verified-2

evidence · `01M6MTKW0005NJ0Z9PPD8RWJMJ` · actor [claude-code-0df2a3a2](../actor/claude-code-0df2a3a2.md)

## result

PASS

## statement

At d01d4252 (attempt 2): go test -trimpath -count=1 ./internal/guide/... passes; build_test.go, review_test.go and plan_test.go assert the `- <<'EOF'` form (plus closing EOF) of every rhiza run findings / rhiza run verdict in build.md, review.md and plan-review.md, and TestGuidesFeedFindingsFromStdin fails on any remaining <file> form. git grep finds no `run findings|verdict ... <file>` left outside .rhiza. The attempt-1 blocker is fixed: go test -trimpath ./cmd/rhiza -run 'TestScripts/guide$' passes with guide.txtar expecting the heredoc forms. Diff read: flags (--leaf, --check, --json) and JSON shapes unchanged.

## Links

### Out

- verifies [guides-stdin-findings](../capability/rhiza/workflows/built-in-workflows/guides-stdin-findings.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
