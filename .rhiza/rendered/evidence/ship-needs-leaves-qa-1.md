# ship-needs-leaves-qa-1

evidence · `01M60B4T59MJRSX8KVPD02382C` · actor [claude-code-f5af7e3e](../actor/claude-code-f5af7e3e.md)

## result

FAIL

## statement

QA cycle 1, binary built from branch head 20f3cee. Fresh `rhiza init` repo with a leafless current change set: `rhiza run start --permission-mode manual` exited 3 with 'has no leaves: ... plan it first' and next 'rhiza run start plan' (--json: code conflict). There was no .git/rhiza/runs and no worktree beyond main. After a capability was added, a run started. The suggested next, `rhiza run start plan`, fails with 'no change set matches "plan"' (qa finding f2). Without --permission-mode, a permission-mode error masks the refusal (f3).

## Links

### Out

- verifies [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
