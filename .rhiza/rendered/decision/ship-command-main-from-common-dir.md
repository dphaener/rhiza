# ship-command-main-from-common-dir

decision · `01M62RFYPDVVVRHSREVFJA6HDY`

## statement

The ship command finds the main checkout from the git common dir, not from the run's worktree folder

## wording

Derive it from the git common dir, the agent's option, picked by the owner. Ship's own target takes the main checkout from the first entry of git worktree list, not from the common dir.

## Links

### Out

- derived-from [ship-step-interview](../source/ship-step-interview.md)

### In

- decided-by [ship-command-after-close-out](../capability/rhiza/ship-command-after-close-out.md)
- supersedes [ship-command-starts-in-main-checkout](ship-command-starts-in-main-checkout.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
