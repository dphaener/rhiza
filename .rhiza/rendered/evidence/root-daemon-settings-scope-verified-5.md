# root-daemon-settings-scope-verified-5

evidence · `01M6SKYHT2E6ES8HVY33N28G29` · actor [claude-code-144aa878](../actor/claude-code-144aa878.md)

## result

PASS

## statement

At 9e1350f3 (attempt 5, fixer commit updating tui-autonomy-app to expect source rhiza.toml): go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-settings-scope|tui-autonomy-app|daemon-settings|daemon-event-hook)$' passes, so the quality blocker (tui-autonomy-app:25 expecting git config) is fixed. daemon-settings-scope sets rhiza.daemon.maxRuns (found in the root's rhiza.toml, absent from git config) and rhiza.autonomy.ship (in project a's git config only), lists both with source and scope, lists rhiza.daemon.trimGoCache as (default, machine), refuses a project key at the root without --project, and with two projects each with its own rhiza.eventCommand runs each hook only for its own project's moves. By hand with a scratch build under a scratch root: daemon set rhiza.autonomy.ship at the root exits 2 with a --project next line; settings --json shows rhiza.daemon.maxRuns scope machine source rhiza.toml, rhiza.autonomy.ship scope project source git config, trimGoCache scope machine. Diff read: one-line expectation change, no code change.

## Links

### Out

- verifies [root-daemon-settings-scope](../capability/rhiza/orchestration/daemon/root-daemon-settings-scope.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
