# docs-reference-config-verified

evidence · `01M6Z8SQ0SQPMKPYPT1NJHVA8P` · actor [claude-code-07585c8d](../actor/claude-code-07585c8d.md)

## result

PASS

## statement

At 23eab5f9 the verifier ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(docs-reference-config|docs-reference-cli|describe-golden)$': all PASS. docs-reference-config.txtar greps out/config.md for '### `rhiza.autonomy.plan`', 'Type: auto or owner' and 'Default: `owner`', and checks that a second run gives the same bytes. go test -trimpath ./internal/api -run TestSettingDefsDescribed passes. The test loops over every settingDefs entry, including the watchdog keys appended in init, and fails any entry with an empty desc or typ. go test ./internal/docs passes. A scratch build's rhiza docs --out wrote config.md with 42 settings (23 base + 10 machine watchdog + 9 project watchdog), each with type, default (unset when empty) and description, plus the --max-runs/--min-free/--permission-mode flag notes. The diff (docs.Config, docsSettings, typ/desc on settingDef and watchdogNames) reads correctly.

## Links

### Out

- verifies [docs-reference-config](../capability/rhiza/documentation/docs-reference-config.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
