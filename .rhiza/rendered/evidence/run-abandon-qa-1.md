# run-abandon-qa-1

evidence · `01M624NVCEK0E1SDPA7GYBP6KQ` · actor [claude-code-7e6e3e8f](../actor/claude-code-7e6e3e8f.md)

## result

PASS

## statement

Built rhiza from the branch head and drove it in a scratch git repo. run abandon --dry-run on a run stopped at a wait printed 'would be abandoned' and left it stopped. run abandon by a 4+ char prefix journalled it failed with detail 'abandoned by <actor>' (run status human and --json); git worktree list and branch cs/<slug> still existed. run resume and run resume --retry on it exited 2 with 'was abandoned'; a second abandon exited 2 'is finished: it failed'. Abandon of a run running in another process exited 3 (running in another process). A SIGKILLed interrupted run was abandoned, after which run start on the change set succeeded and finished.

## Links

### Out

- verifies [run-abandon](../capability/rhiza/workflows/workflow-runs/run-abandon.md)

## History

- `01M5ZNV4E4AJR52EX631180DMS` A change set has at most one live run: rhiza run start refuses a second run on a change set whose run is still going, naming that run.
