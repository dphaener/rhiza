# watchdog-policy-qa-1

evidence · `01M6WZEQXRV9R55RZEGZKYP42W` · actor [claude-code-9eef2845](../actor/claude-code-9eef2845.md)

## result

PASS

## statement

Built rhiza at the branch head in a scratch root and project. rhiza daemon settings listed all ten rhiza.daemon.watchdog.* machine keys and nine rhiza.project.watchdog.* project keys (no pollBehind) with their defaults. daemon set rhiza.daemon.watchdog.agentNudge 7m wrote [daemon] "rhiza.daemon.watchdog.agentNudge" = "7m" to rhiza.toml. A hand-edited value of "soon" made both daemon settings and daemon --once fail with an error naming the key. daemon set ... commandLimit 0s was rejected with the key named. Setting rhiza.project.watchdog.enabled false in git config showed (git config, project), and the daemon then left an autonomous stalled run alone. A workflow step with stall = "5s" was used as that step's limit.

## Links

### Out

- verifies [watchdog-policy](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-policy.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
