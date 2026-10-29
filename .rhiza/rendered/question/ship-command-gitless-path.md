# ship-command-gitless-path

question · `01M62RM4HA219N0SDXD961BBPG`

## statement

From a worktree folder that close-out has removed, git cannot run, so how does the ship command find the main checkout's git common dir?

## Links

### Out

- blocks [ship-command-after-close-out](../capability/rhiza/ship-command-after-close-out.md)

### In

- answers [ship-command-starts-in-main-checkout](../decision/ship-command-starts-in-main-checkout.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
