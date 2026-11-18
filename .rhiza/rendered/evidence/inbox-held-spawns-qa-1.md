# inbox-held-spawns-qa-1

evidence · `01M7QK1V52V4KW3VPMGAT6MKWC` · actor [claude-code-062a7236](../actor/claude-code-062a7236.md)

## result

PASS

## statement

QA with a scratch build: held step sessions and a held verifier spawn of a live run each showed one spawn-held item in rhiza inbox and --json with change set, step, role and the hold's reading and limit, stable ID; gone after release and run end. Minor: rhiza inbox open on the item says no inbox item

## Links

### Out

- verifies [inbox-held-spawns](../capability/rhiza/orchestration/inbox/inbox-list/inbox-held-spawns.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
