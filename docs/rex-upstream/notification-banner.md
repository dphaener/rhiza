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
