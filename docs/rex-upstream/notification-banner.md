# Rex upstream request: reliable macOS banners for blocked reports

Status: to file with Rex upstream. Owner decision:
`vm-move-notify-chime-and-marker` (Rhiza change set "Move this repo's
Rhiza", leaf `rhiza-on-vm`).

## Context

Rhiza's daemon and agents run on a VM. The owner works from the Mac's Rex
app connected to the VM's Rex host. When an agent in a VM tab raises a
question it reports `blocked` through Rex's program-status protocol
(OSC 7501, rev 0.4).

## What Rex does today

Rex's program-status protocol leaves presentation to the terminal and
promises no macOS banner. In the owner's live test on 2026-11-05:

- every blocked report from a VM tab chimed on the Mac when the tab was not
  focused (tests C to G);
- one showed a macOS banner (test B).

The blocked marker on the tab and in the session picker is always shown.
For this milestone the chime plus the marker is accepted as the
notification.

## Request

Deliver a macOS notification banner reliably, not only sometimes, when a
program in a tab that is not focused reports `blocked` (or another state
that needs the user), including for tabs hosted on a remote Rex server the
Mac app is connected to. Document in the protocol whether a banner is
promised, and offer a way to turn it on or off per session.

## Why

Rhiza's autonomous runs ask the owner questions. A banner reaches the owner
when the Rex window is hidden or on another Space; a chime alone can be
missed.

---

# Background: the investigation notes

The notes below were gathered while testing banners from the VM (Nov 3 to 5). They are kept here as background for the request above.

## Rex: blocked status from a remote (VM) host chimes but posts no banner

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
(`ESC` = 0x1b; the report ends with ST, `ESC \` = 0x1b 0x5c, as `rexRecord`
in `internal/api/run.go` always writes it):

```
ESC ]7501;state=blocked:kind=permission:app=rhiza:title=cmVwbzogQWdlbnQgcnVucyDCtyBydW4gcnVuLTEgwrcgYnVpbGQ=:msg=QWdlbnQgcnVuczogYnVpbGQgd2FpdHMgb24geW91IGluIHRhYiByZXBvOiBBZ2VudCBydW5zIMK3IHJ1biBydW4tMSDCtyBidWlsZDogQ2xhdWRlIG5lZWRzIHlvdXIgcGVybWlzc2lvbiB0byB1c2UgQmFzaA==ESC \
```

Decoded: title = `repo: Agent runs · run run-1 · build`; msg =
`Agent runs: build waits on you in tab repo: Agent runs · run run-1 · build: Claude needs your permission to use Bash`.
A question block is the same with `kind=question`.

### Claude Code's own blocked report

Captured on the VM on 2026-11-04 with Claude Code 2.1.316: a scratch Rex session
on the VM ran interactive `claude --permission-mode default "Run the bash
command: touch captured.txt"` while `rex events -s <session> --json` streamed
its events. Claude Code reported idle, then working, then this line when its
Bash permission prompt came up:

```
{"type":"session_event","payload":{"session_id":"session:01a1a568-f1d5-79f4-8ce7-dcc1a3815925","event":{"type":"block_event","payload":{"block_id":"block:01a1a568-f1d5-7bef-89b6-e7063ca1311d","creator_name":"com.superlogical.terminal","name":"program_status_changed","payload":{"id":"","record":{"state":"blocked","kind":"permission","app":"claude-code","effective_app":"claude-code","msg":"approve Bash: touch captured.txt","updated_at":"2026-11-04T07:16:43.632990895Z","owner":{"pid":1329432,"name":"claude"}},"reason":"report"}}}}}
```

Note: Claude Code's record carries `state`, `kind` and `msg` but no `title`;
rhiza's report carries a `title`. If Rex needs a title to post a banner, that
alone may explain a chime with no banner for Claude Code's own report.

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

## iPhone: no notifications at all (2026-11-05)

A titled blocked report (as rhiza writes it) and an untitled one (as Claude
Code writes it) were sent a minute apart from a scratch tab on the VM while
the owner watched an iPhone 15 Pro running Rex Dev 0.1.0 (1779) with the app
in the background. Neither produced a notification. `rex client list` showed
the iPhone connected for under a minute (the Mac's Rex Beta 0.1.0 (1758) had
held its connection for 7 hours), so the phone app appears to reconnect only
when opened. iOS Settings has no notification entry for Rex, so the app has
never asked for notification permission. Request: notifications on iOS for
blocked programs, which while the app is suspended means push (APNs) from
the server.

## iPhone: no floating layers (2026-11-05)

On the same iPhone build (Rex Dev 0.1.0 (1779)), apps that open as Rex
floating layers don't work. Rhiza's five apps (inbox, request, runs, daemon
and autonomy) each open as a rexuikit `surface` Modal over the current pane.
On the phone there is no modal: the app's content renders briefly, and
control goes straight back to the terminal. Owner, 2026-11-05: "Things kind
of render, but there aren't modals, control is just handed right back to the
terminal." The phone app also caches the server's action list until it is
reopened, so new actions don't appear after a reload. Request: support
floating layers (Modal, Drawer, Split) in the iOS app, or report that they
are unsupported, so a program can fall back to drawing in its own pane.
