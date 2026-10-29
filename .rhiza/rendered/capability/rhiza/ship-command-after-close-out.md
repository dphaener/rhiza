# ship-command-after-close-out

capability · `01M62RFZ5YEP4PWEC0799MK0WZ` · parent [rhiza](../rhiza.md)

## acceptance

A testscript removes the worktree, runs rhiza ship command from the vanished folder and sees the script run in the main checkout with exit 0; a ship-change run whose ship step follows close-out ends succeeded

## statement

After close-out has removed a change set's worktree, rhiza ship command still runs the repo's ship script, in the main checkout at the same project folder

## Links

### Out

- decided-by [ship-command-main-from-common-dir](../../decision/ship-command-main-from-common-dir.md)
- decided-by [ship-command-starts-in-main-checkout](../../decision/ship-command-starts-in-main-checkout.md)

### In

- blocks [ship-command-gitless-path](../../question/ship-command-gitless-path.md)
- verifies [ship-command-after-close-out-verified](../../evidence/ship-command-after-close-out-verified.md)
- verifies [ship-command-after-close-out-qa-1](../../evidence/ship-command-after-close-out-qa-1.md)
- verifies [ship-command-after-close-out-verified-2](../../evidence/ship-command-after-close-out-verified-2.md)
- verifies [pr-73-ship-step-after-close-out](../../evidence/pr-73-ship-step-after-close-out.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
