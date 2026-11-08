# warn-dont-safeguard

decision · `01M6W4H6F2KYTFE6ER04FQ3EWE`

## options

```json
[
  "warn, never block",
  "fence untrusted request text",
  "refuse risky setups"
]
```

## quote

"i don't think we can reliably safeguard against users shooting themselves in the foot. best we can do is warn people." (owner, 2026-11-07); "a schedule can run anything. who cares?" (owner, 2026-11-07). The agent's proposed defaults (warnings never block; schedules per repo in rhiza.toml, run once after downtime, skip a slot while the previous run is going) were accepted by the owner ("go ahead").

## statement

Rhiza does not try to safeguard owners against risky setups (such as an autonomous workflow with bypassPermissions fed by outside input); it warns plainly, in rhiza health, workflow check and the guides, and leaves the choice to the owner. A schedule can run anything: a workflow or a plain command.

## Links

### In

- decided-by [workflows](../capability/rhiza/workflows.md)
- decided-by [watchdog-policy](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-policy.md)

## History

- `01M6W4H2QV52XMGY4JXMJ11A3Q` Record the owner's rule: Rhiza warns about risky setups instead of safeguarding, and a schedule can run anything
