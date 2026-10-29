# ship-command-after-close-out-qa-1

evidence · `01M6427NVDNXAX1HH2J0NSRS2W` · actor [claude-code-edd7b77f](../actor/claude-code-edd7b77f.md)

## result

PASS

## statement

QA by hand at e47d48c with a built binary (go build ./cmd/rhiza), outside testscript: in a scratch sandbox (bare origin, fake rex/claude/gh from run-ship-change.txtar, RHIZA_* unset) ran rhiza run start on a two-leaf accepted change set, then rhiza run approve from inside the change set's worktree repo.ship-it. Close-out removed repo.ship-it, then the ship step ran .rhiza/workflows/ship.sh in the main checkout: ship.log holds 'shipped' once, order.log 'main-green, main-green, shipped', run status shows 'close-out: pass, then ship: exit 0' and 'ship: pass, then end: exit 0', and the run state is succeeded. The same steps with a binary built from the parent commit 1a60493 fail: 'ship: error ... chdir .../repo.ship-it: no such file or directory', run state failed, no ship.log.

## Links

### Out

- verifies [ship-command-after-close-out](../capability/rhiza/ship-command-after-close-out.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
