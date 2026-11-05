# guides-stdin-findings-qa-1

evidence · `01M6MTE66PQX4XQHCCZPG1NZHR` · actor [claude-code-a64abfd3](../actor/claude-code-a64abfd3.md)

## result

PASS

## statement

At 6521b137 built ./cmd/rhiza (go build -trimpath) and ran its guide build, guide review and guide plan-review: all five rhiza run findings / rhiza run verdict commands show the `- <<'EOF'` ... `EOF` heredoc form with their flags (--leaf, --check, --json) and unchanged JSON shapes; no `<file>` form of either remains. Hands-on, ran the review guide's exact command `rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'` with the built binary in this qa session: it recorded round 1 of check qa (exit 0). The plan-review form `rhiza run findings - --json <<'EOF'` parses (reaches the RHIZA_RUN check) in a scratch repo. go test -trimpath -count=1 ./internal/guide/... ok.

## Links

### Out

- verifies [guides-stdin-findings](../capability/rhiza/workflows/built-in-workflows/guides-stdin-findings.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
