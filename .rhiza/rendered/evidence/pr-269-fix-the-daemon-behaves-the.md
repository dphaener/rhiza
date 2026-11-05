# pr-269-fix-the-daemon-behaves-the

evidence · `01M6MKYYF493C64HBQJC1QTJVB` · actor [dphaener](../actor/dphaener.md)

## head

bcb75e8e403c6a91f61513244da0ba829851d9e8

## pr

https://github.com/dphaener/rhiza/pull/269

## result

1 check passed: check SUCCESS

## statement

Pull request #269 (https://github.com/dphaener/rhiza/pull/269) passed CI at bcb75e8

## status

verified

## Links

### Out

- verifies [daemon-service-env](../capability/rhiza/orchestration/daemon/daemon-service-env.md)

## History

- `01M6JY2NZESHV22Y931D7XMBEE` Fix The daemon behaves the same however its installation runs it: the systemd unit and the Rex tab start rhiza daemon with the same arguments, so choosing service or tab never changes whether runs go on without the owner (replaces draft change set 01M6FN21KA287N, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
