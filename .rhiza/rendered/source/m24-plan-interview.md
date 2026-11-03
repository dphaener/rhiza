# m24-plan-interview

source · `01M6FEF610DN2BRHAMYJX9BSDG`

## quotes

```json
[
  {
    "by": "owner",
    "text": "It shouldn't run inside a tab",
    "where": "answer on where run tabs open when the daemon has no tab"
  },
  {
    "by": "owner",
    "text": "The user can choose. Is this a server only installation, or do I want to just run the daemon in a tab on my mac (when Rex for Linux is available - the UI - we can expand to that as well). It needs to be flexible.",
    "where": "answer on what starts and stops the daemon"
  },
  {
    "by": "owner",
    "text": "#1 is fine for now but I think the future should be it checks for an actual new release. Since we don't have that yet we can just check main.",
    "where": "answer on who updates the VM's rhiza"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Headless daemon + VM move: rhiza daemon runs as a systemd service on the VM, opens run tabs in a named session on the VM's Rex server, and the Mac's Rex app works it all as a remote host"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "rhiza daemon install writes a systemd user unit for this repo, plus rhiza daemon uninstall"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Fixed name per repo: always 'rhiza-<repo name>', no setting"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Rex only: rely on Rex for notifications; if they don't come through, raise it with Rex upstream"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "New Ubuntu 26 VM on Proxmox, reached over Tailscale, Rex server from the Ubuntu package with rex.service enabled"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "HTML runbook like afk's: a standalone HTML runbook in ~/code beside the others, and Evidence in the model"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Plan all three change sets now: daemon as a service, self-update from main, move to the VM"
  }
]
```

## raw_sha256

bd5ecc930d30bf776ab487e21655ddcb654301ff7cf4895d82ba64b04a928dbd

## session

841f57a3-6242-4b77-8c1f-01a19edced79

## summary

Planned the milestone 'Make Rex run as a background daemon on Linux and move this installation to a dedicated VM'. Rex's Linux server already runs under systemd and keeps sessions across restarts, so the work is Rhiza's: the daemon must run without a Rex tab of its own. Decided: done means a headless daemon plus actually moving this installation; on a server the daemon is a systemd service installed by rhiza daemon install, never in a tab, but each installation chooses service or tab (the Mac keeps the tab), so start/stop must handle both; run tabs open in a fixed per-repo Rex session; the daemon self-updates from main for now, from releases later; notifications stay Rex only; a new Ubuntu Server 26 VM on Proxmox; the move is captured in an HTML runbook like the afk one. Three change sets, all planned now: daemon as a service, self-update from main, move to the VM. Open: who carries out the VM build itself.

## title

Milestone 24 plan: Rhiza's daemon headless on a Linux VM, with the owner

## Links

### In

- derived-from [m24-server-daemon-no-tab](../decision/m24-server-daemon-no-tab.md)
- derived-from [m24-service-or-tab](../decision/m24-service-or-tab.md)
- derived-from [m24-update-from-main](../decision/m24-update-from-main.md)
- derived-from [answer-841f57a3-1](../decision/answer-841f57a3-1.md)
- derived-from [answer-841f57a3-2](../decision/answer-841f57a3-2.md)
- derived-from [answer-841f57a3-3](../decision/answer-841f57a3-3.md)
- derived-from [answer-841f57a3-4](../decision/answer-841f57a3-4.md)
- derived-from [answer-841f57a3-5](../decision/answer-841f57a3-5.md)
- derived-from [answer-841f57a3-6](../decision/answer-841f57a3-6.md)
- derived-from [answer-841f57a3-7](../decision/answer-841f57a3-7.md)
- derived-from [answer-841f57a3-8](../decision/answer-841f57a3-8.md)
- derived-from [answer-841f57a3-9](../decision/answer-841f57a3-9.md)
- derived-from [answer-841f57a3-10](../decision/answer-841f57a3-10.md)
- derived-from [answer-841f57a3-11](../decision/answer-841f57a3-11.md)
- derived-from [answer-841f57a3-12](../decision/answer-841f57a3-12.md)
- derived-from [answer-841f57a3-13](../decision/answer-841f57a3-13.md)

## History

- `01M6FDSQH550AN63HK4MK00T6V` Make Rex run as a background daemon on Linux and move this installation to a dedicated VM
