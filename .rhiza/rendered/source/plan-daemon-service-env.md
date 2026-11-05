# plan-daemon-service-env

source · `01M6JY9CRQHY97RAHDCSHRS9V1`

## quotes

```json
[]
```

## raw_sha256

fc8cd946077b7cc3621202d13f00886458cce7172502a8ac303463d11bb40240

## session

db198e1f-2e16-40ae-8f4b-01a1a5e18368

## summary

Autonomous plan run, so no owner took part. The planner found that the systemd unit's ExecStart and the Rex tab already start rhiza with the same daemonArgs (daemon --autonomous) since #200, and both scripts check it. The one real difference between service and tab is the environment: daemon install writes no PATH or TMPDIR, and the VM works only because of a hand-written drop-in. The planner asked the arbiter for the scope. The arbiter ruled (Decision daemon-service-env-matches-tab, assumed) that the leaves cover the environment: daemon install writes the installing shell's PATH and TMPDIR into the unit as Environment= lines. The shared args get no leaf of their own. The change set serves milestone-25, relates to the leafless draft 01M6FN21KA287N, and that draft is withdrawn.

## title

Plan of 'The daemon behaves the same' (autonomous, arbiter)

## Links

### In

- derived-from [daemon-service-env-matches-tab](../decision/daemon-service-env-matches-tab.md)

## History

- `01M6JY2NZESHV22Y931D7XMBEE` Fix The daemon behaves the same however its installation runs it: the systemd unit and the Rex tab start rhiza daemon with the same arguments, so choosing service or tab never changes whether runs go on without the owner (replaces draft change set 01M6FN21KA287N, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
