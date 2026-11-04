# Rex: blocked status from a remote (VM) host chimes but posts no banner

DRAFT for Rex upstream. The VM-side facts below were captured on the VM; the
**OWNER** fields are blank for the owner to fill on the Mac. File this with Rex
only if, after the Mac checks, a question in a VM agent tab still only chimes.

## Symptom

A question asked in a Claude Code agent tab on the VM plays the notification
chime on the Mac, but no macOS notification banner appears. Owner, 2026-11-03:
"I have seen no notification on my mac. But I have heard chimes".

## Setup

- Rex server on the VM (`rex version`):
  ```
  rex 0.0.1793645657-gd018d2
  commit: d018d2de8cc31281265906a4648ec9cae8d2f4c5
  built: 2026-11-02T18:54:17Z
  api: revision 10 (supports 0 to 10)
  ```
- VM: host `rhiza`, Ubuntu 26.04, Rex server account `rhiza`, reached over
  Tailscale (`rex.tail721b54.ts.net`, HTTPS available, tailnet members not
  allowed). Agent tabs run Claude Code under a `rhiza session run` wrapper.
- The Mac runs the Rex app attached to the VM as a remote host (`rex -S rhiza`).
  The Mac runs no daemon and no agents.

## What the VM writes

### rhiza's blocked report (root record, app=rhiza)

From the built binary of this branch (testscript `session-run-blocked-notice`,
head be648824), for a permission block in run `run-1`, step `build`. Exact bytes
(`ESC` = 0x1b, terminated by BEL or ST as the writer emits):

```
ESC ]7501;state=blocked:kind=permission:app=rhiza:title=cmVwbzogQWdlbnQgcnVucyDCtyBydW4gcnVuLTEgwrcgYnVpbGQ=:msg=QWdlbnQgcnVuczogYnVpbGQgd2FpdHMgb24geW91IGluIHRhYiByZXBvOiBBZ2VudCBydW5zIMK3IHJ1biBydW4tMSDCtyBidWlsZDogQ2xhdWRlIG5lZWRzIHlvdXIgcGVybWlzc2lvbiB0byB1c2UgQmFzaA==
```

Decoded: title = `repo: Agent runs · run run-1 · build`; msg =
`Agent runs: build waits on you in tab repo: Agent runs · run run-1 · build: Claude needs your permission to use Bash`.
A question block is the same with `kind=question`.

### Claude Code's own blocked report

NOT CAPTURED on the VM: a blocked state needs a live interactive Claude Code
permission or question prompt, which a headless build run cannot raise without
opening a real agent window. To capture it, ask a question in a VM agent tab
and run on the VM:

```
rex events -s <session> --json | grep '"state":"blocked"'
```

Rex delivers it nested as
`{"type":"session_event","payload":{"session_id":...,"event":{"type":"block_event","payload":{"block_id":"block:...","creator_name":"com.superlogical.terminal","name":"program_status_changed","payload":{"id":"","record":{"state":"blocked","kind":...,"msg":...,"updated_at":...},"reason":"report"}}}}}`
(the shape rhiza's own reports come back in, captured live earlier). Paste the
record for Claude Code's blocked report here:

```
<PASTE rex events --json LINE>
```

## Mac observations (OWNER fills in)

- Mac Rex app version: ______
- Does a banner now show for a question in a VM agent tab (after rhiza's
  agent-tab-blocked-notice change is installed on the VM)? yes / no
- Rex window focused when the question was asked? yes / no
- macOS notification permission for Rex (System Settings > Notifications):
  allowed / banners / alerts / none
- Do Focus or Do Not Disturb modes apply? ______
- Does a blocked report from a local (Mac) block get a banner? yes / no

## Ask of Rex

If a remote host's block reports `state=blocked` and the Mac plays the chime,
post the macOS banner too, or document the rule that suppresses it for remote
hosts or focused windows.
