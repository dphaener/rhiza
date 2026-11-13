# guides-stdin-findings

capability · `01M6MMAF8Z4XWD1JD9WG5PRYBR` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/guide/... passes: TestGuidesFeedFindingsFromStdin asserts the verifier, check-compliance, critic and plan-review topics show the `<<'EOF'` stdin form and no `<file>` form of rhiza run findings or rhiza run verdict; TestCheckComplianceText, TestCheckQualityText, TestCheckQAText and TestCheckModelText assert each check topic shows `rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`, TestCriticText that the critic topic shows `rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'`, and TestReviewRoleText finds both heredocs and {"findings":[]} in the role topics

## statement

The verifier guide topic, the plan-review guide topic, the four check topics (check-compliance, check-quality, check-qa, check-model) and the critic topic show rhiza run findings and rhiza run verdict fed from stdin by a heredoc (`- <<'EOF'` ... `EOF`), with their existing flags (--leaf, --check, --json) and JSON shapes ({"findings":[]} when clean), so an agent records findings and verdicts without a temp file or compound shell command; the review index shows neither

## Links

### Out

- decided-by [stdin-guides-ruling](../../../../decision/stdin-guides-ruling.md)

### In

- blocks [review-01m79a41352xbz4fc2fjt5903j-r3-f6](../../../../question/review-01m79a41352xbz4fc2fjt5903j-r3-f6.md)
- verifies [guides-stdin-findings-verified](../../../../evidence/guides-stdin-findings-verified.md)
- verifies [guides-stdin-findings-qa-1](../../../../evidence/guides-stdin-findings-qa-1.md)
- verifies [guides-stdin-findings-verified-2](../../../../evidence/guides-stdin-findings-verified-2.md)
- verifies [pr-282-fix-agents-in-runs-hit-claude](../../../../evidence/pr-282-fix-agents-in-runs-hit-claude.md)
- verifies [pr-521-build-role-topics](../../../../evidence/pr-521-build-role-topics.md)
- verifies [pr-537-review-role-topics](../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
