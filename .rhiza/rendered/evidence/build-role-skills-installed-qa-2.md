# build-role-skills-installed-qa-2

evidence · `01M79C8N5ZCK3MZE867ZWX7JS7` · actor [claude-code-9b656f6b](../actor/claude-code-9b656f6b.md)

## result

PASS

## statement

With the built binary in a scratch repo (HOME isolated), `rhiza guide install --agent claude` wrote skills/rhiza-builder, rhiza-verifier, rhiza-build-retry and rhiza-build SKILL.md (10 written); re-running reported 0 written, 10 unchanged. `rhiza guide build` is a short index naming /rhiza-builder, /rhiza-build-retry and /rhiza-verifier and their rhiza guide routes. go test ./internal/api -run GuideInstall ok; TestScripts/(guide|guide-install|build-evidence-commit) ok

## Links

### Out

- verifies [build-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-role-skills-installed.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
