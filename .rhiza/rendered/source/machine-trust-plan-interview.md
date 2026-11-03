# machine-trust-plan-interview

source · `01M6G04YCB7DCPA75XAMJB0DKE`

## quotes

```json
[
  {
    "by": "orchestrator (owner's proxy, not the owner)",
    "text": "the owner is asleep; these answers come from the orchestrator acting as the owner's proxy under their standing instruction (\"act on my behalf\"). Record each as the orchestrator's (arbiter) wording, assumed for the owner, never as owner quotes. Owner intent tonight: this VM is a dedicated machine and every agent session on it should run in bypassPermissions.",
    "where": "plan session, mid-interview message"
  },
  {
    "by": "agent's options, picked by the orchestrator for the owner",
    "text": "milestone-24; ~/.config/rhiza/machine.toml; hand-edit only; missing = daily, bad = error; no default change, only lifts the refusal; rhiza machine show and a health part; read live at each check",
    "where": "plan session, two AskUserQuestion rounds"
  }
]
```

## raw_sha256

84706732b043905a13842aeaf1a1b57f8123a709b0ec8b287a2bbdb08c8ef6dc

## session

5a563475-f1ff-4b44-9c5f-01a19fe520aa

## summary

Planned change set Machine trust. The owner was asleep: an orchestrator answered as the owner's proxy under the owner's standing instruction to act on their behalf, so every pick below is the agent's option assumed for the owner, never the owner's words. Settled: the change set serves milestone-24. A machine's trust lives in a machine config file, $XDG_CONFIG_HOME/rhiza/machine.toml (default ~/.config), as trust = "dedicated" | "daily". It is set only by hand-editing that file: rhiza has no command or API operation that writes trust. A missing file means daily; a malformed file or unknown value is an error naming the file, and bypass stays refused. Trust changes no default permission mode, it only lifts the refusal of bypassPermissions. Trust shows in a read-only rhiza machine show and in a machine part of rhiza health. Trust is read live at every check (session start, session run, run start, daemon start), and no session record carries it. The orchestrator said the owner's intent tonight is that this VM is dedicated and every agent session on it runs in bypassPermissions: that is done on the VM by writing the machine config and setting rhiza.daemon.permissionMode=bypassPermissions after this lands, since no built-in workflow step sets its own mode. Nothing left open.

## title

Machine trust plan, with the orchestrator as the owner's proxy

## Links

### In

- derived-from [mt-serves-m24](../decision/mt-serves-m24.md)
- derived-from [machine-config-file](../decision/machine-config-file.md)
- derived-from [trust-hand-edited](../decision/trust-hand-edited.md)
- derived-from [trust-missing-or-bad](../decision/trust-missing-or-bad.md)
- derived-from [trust-lifts-refusal-only](../decision/trust-lifts-refusal-only.md)
- derived-from [trust-shown](../decision/trust-shown.md)
- derived-from [trust-read-live](../decision/trust-read-live.md)
- derived-from [machine-config](../capability/rhiza/runtime/machine-config.md)
- derived-from [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)
- derived-from [machine-show](../capability/rhiza/runtime/machine-show.md)
- derived-from [health-machine-part](../capability/rhiza/runtime/health-machine-part.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
