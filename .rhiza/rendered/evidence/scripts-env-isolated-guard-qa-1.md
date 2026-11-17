# scripts-env-isolated-guard-qa-1

evidence · `01M7M4A2FYT539TBWTAATZQJ89` · actor [claude-code-22f52dfa](../actor/claude-code-22f52dfa.md)

## result

PASS

## statement

QA at 53f29ac4f, session RHIZA_*/SHIP_*/REX_* unset, under a capped scope with -p 2 -trimpath: go test ./cmd/rhiza -run '^TestScriptEnvIsolated$' passes. With keepScriptEnv temporarily also keeping every os.Environ() name, it fails 'scriptEnv: leaked from the test process: RHIZA_CHANGE (in RHIZA_CHANGE=), ... REX_BLOCK (in REX_BLOCK=)', naming all eight leakVars (reverted; tree clean). -run 'TestScripts/^session-run$' passes. Edge cases: -count=2 passes twice; '^TestScriptEnvIsolated$|TestScripts/^session-run$' in one process passes both; run with the session's real RHIZA_*/REX_*/SHIP_* still set it passes and the logged script env holds none of them.

## Links

### Out

- verifies [scripts-env-isolated-guard](../capability/rhiza/integrity/self-check/scripts-env-isolated-guard.md)

## History

- `01M7M2YWF4KEF5RV32ZS0C4XK7` Fix cmd/rhiza testscripts inherit the calling session's RHIZA_*, SHIP_* and REX_* variables, so a script run from inside a run session fails or passes differently than in CI
