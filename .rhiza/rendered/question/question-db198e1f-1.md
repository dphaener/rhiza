# question-db198e1f-1

question · `01M6JY66RRW4DKBTT4X2AFA68A`

## session

db198e1f-2e16-40ae-8f4b-01a1a5e18368

## statement

Change set 'The daemon behaves the same' (01M6JY2): the code already does what its title says. internal/api/daemon_startstop.go has one daemonArgs = [daemon --autonomous], and both the systemd unit's ExecStart and the Rex tab use it, since #200. daemon-install.txtar and daemon-start-stop.txtar each grep for 'daemon --autonomous'. The one real difference between service and tab is the environment. On the VM, PATH and TMPDIR come from a drop-in written by hand (~/.config/systemd/user/rhiza-rhiza.service.d/env.conf), and daemon install doesn't write it. So a fresh install would run agents with systemd's bare PATH, unlike the tab. What should this change set's leaves cover?

## Links

### In

- answers [daemon-service-env-matches-tab](../decision/daemon-service-env-matches-tab.md)

## History

- `01M6JY2NZESHV22Y931D7XMBEE` Fix The daemon behaves the same however its installation runs it: the systemd unit and the Rex tab start rhiza daemon with the same arguments, so choosing service or tab never changes whether runs go on without the owner (replaces draft change set 01M6FN21KA287N, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
