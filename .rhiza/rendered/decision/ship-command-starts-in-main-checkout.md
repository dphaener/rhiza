# ship-command-starts-in-main-checkout

decision · `01M640VS2H4XZ8Z0SABYJT53ZN`

## reasoning

git cannot run from a folder that is gone, so the answer has to be known before close-out. The run already depends on the main checkout (its journal and ship state live in its .git/rhiza), and the change set's intent already says steps after close-out start in the main checkout. Retries started from the main checkout today (#66, #67, #68) passed the ship step, which is this behaviour done by hand.

## statement

The run finds the main checkout while its worktree still exists: the runner records the main checkout (the folder holding the git common dir, where the run's journal lives under .git/rhiza) when the run starts, and starts every step after close-out there. The ship command then runs in a folder where git works and finds the common dir as it does from any checkout; it never has to start git from the removed folder.

## wording

decided by the orchestrator acting as arbiter, under the owner's standing word to run autonomously: "I'm stopping for the night, let's just run autonomously from here on out, I'm just accepting the recommended answers anyway"

## Links

### Out

- answers [ship-command-gitless-path](../question/ship-command-gitless-path.md)
- supersedes [ship-command-main-from-common-dir](ship-command-main-from-common-dir.md)

### In

- decided-by [ship-command-after-close-out](../capability/rhiza/ship-command-after-close-out.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
