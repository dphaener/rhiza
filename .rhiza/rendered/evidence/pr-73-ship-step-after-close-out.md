# pr-73-ship-step-after-close-out

evidence · `01M647BXK8065T7G5SSB897DH6` · actor [darin-haener](../actor/darin-haener.md)

## head

cb383d210a226bd705faf6b0ac6d1f97108bc175

## pr

https://github.com/dphaener/rhiza/pull/73

## result

1 check passed: check SUCCESS

## statement

Pull request #73 (https://github.com/dphaener/rhiza/pull/73) passed CI at cb383d2

## status

verified

## Links

### Out

- verifies [ship-command-after-close-out](../capability/rhiza/ship-command-after-close-out.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
