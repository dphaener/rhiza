# inbox-held-spawns-qa-2

evidence · `01M7QKYRHJ78B6RYVY7TE4WQ5Q` · actor [claude-code-608090ea](../actor/claude-code-608090ea.md)

## result

PASS

## statement

QA cycle 2, scratch build: a live run's held verifier spawn showed one spawn-held item in rhiza inbox and --json naming change set, title, step, parent session, role and 'machine hold: memory-low: 1.0 GiB available, under the minimum of 2.0 GiB'; rhiza inbox open on it (f1 fix 13a633dfe) focused the run's journalled tab block:qa1 (opened true) when the run had one, and without a tab was refused as usage naming the held session and run; a bogus ID says no inbox item; the item was gone after the spawn started and after the run ended

## Links

### Out

- verifies [inbox-held-spawns](../capability/rhiza/orchestration/inbox/inbox-list/inbox-held-spawns.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
