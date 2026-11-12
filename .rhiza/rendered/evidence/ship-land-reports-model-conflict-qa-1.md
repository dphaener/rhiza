# ship-land-reports-model-conflict-qa-1

evidence · `01M77SJEC2BVP7TVKB3MKKAS51` · actor [claude-code-e1432792](../actor/claude-code-e1432792.md)

## result

PASS

## statement

QA by hand with the binary built at 668f5cf08, in a scratch repo (bare origin, fake gh reporting PR #7 green): change set Ship it set first's statement, ship check passed, then change set Other rewrote first's statement and landed. On branch cs/ship-it, rhiza ship land exited 1 with 'Refused to land change set ...: its own operations rest on a stale base: <node> (field:statement) rewritten by op <op> of change set <Other>'; --json had landed=false and one model_conflicts entry (op, node, fields, base, landed_op, landed_change); the ship state file recorded the same; nothing was committed and gh pr merge was never called. rhiza ship check --model-conflict then exited 0 (JSON detail names the conflict); in another repo, after rhiza set first --rebase it exited 1 ('the recorded model conflict is repaired'), and with no conflict recorded it exited 1. Minor: its human line with a PR recorded reads 'PR #7 at : open from cs/ship-it' instead of naming the conflict.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
