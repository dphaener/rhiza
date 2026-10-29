# ship-step-interview

source · `01M62RFR9HBFWXACG8PSYBNYNX`

## quotes

```json
[]
```

## raw_sha256

f2a3cb66ebd9cc7fc3bc617862e0ccb1e6e4215095d649df49611809ef6cffd1

## session

0843113a-2a52-485f-8da5-01a185854c07

## summary

The ship-change run's last step, ship, fails after close-out because `rhiza ship command` starts from the run's worktree, which close-out has removed. shipCommandDir runs git from that folder, so it falls back to the gone path. Owner decision: the ship command finds the main checkout from the git common dir, and keeps the project-root subfolder. Ship's own target does not work that way: it takes the main checkout from the first entry of git worktree list and uses the common dir only for the state path. The owner chose this option from the agent's offered options (accepted by the owner); no owner words were typed. Still open: nothing.

## title

Ship step after close-out, with the owner

## Links

### In

- derived-from [ship-command-main-from-common-dir](../decision/ship-command-main-from-common-dir.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
