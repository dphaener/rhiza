# deploy-check-daemon-setting

decision · `01M6J31W2SVSV4GWA64TZR8JPE`

## statement

The ship deploy check reads its health URL and timeout from per-repo daemon settings; unset means no deploy check.

## wording

Daemon setting: per-repo settings via rhiza daemon set (git config), for example ship.health-url and ship.deploy-timeout; unset means no deploy check. The assistant's option, picked by the owner.

## Links

### Out

- derived-from [milestone-26-interview](../source/milestone-26-interview.md)

### In

- decided-by [ship-deploy-settings](../capability/rhiza/orchestration/shipping/ship-deploy-settings.md)
- decided-by [ship-deploy-check](../capability/rhiza/orchestration/shipping/ship-deploy-check.md)

## History

- `01M6J2RTBPHTT1VFJYW2ZA5FRC` Milestone 26: an existing repo with built code comes into Rhiza, Nebari first, and its work then runs on Rhiza workflows
