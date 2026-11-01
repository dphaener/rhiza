# m17-orchestrator-picks

source · `01M6BV24TK47JGT57WK281CG1J`

## quotes

```json
[
  {
    "by": "owner",
    "text": "if you run out of work keep picking high value tasks. don't stop and wait for my input.",
    "where": "orchestrator session, 2026-10-31"
  }
]
```

## raw_sha256

b22a785523a97bdc0e72231e9c8199180364aa9cdb09d7edeaa96baedd840dc3

## summary

With milestone 16 the last one queued, the orchestrator planned milestone 17 under the owner's standing instruction. The founding conversation decided there are hooks on model events from the start ('yes to all four', B6); they sit in the Area hooks umbrella, unbuilt. Milestone 14's notify hook already runs an owner-configured command per new inbox item from the daemon. Every lifecycle transition (propose, accept, reject, withdraw, land, reopen) is written in one place, but firing hooks there would only fire on the writer's machine and inside agent worktrees, so the daemon walks the folded lifecycle each poll instead, which also catches entries that arrive by git pull. Entry ULIDs can be older than a cursor when pulled from another branch, so the daemon remembers which lifecycle entries it has handled, seeding them all on first run so history is not replayed. model-event-hooks cannot be carved without its parent hooks, which exists only in the umbrella, so both move; drift detection is not built, so drift events stay with Area integrity.

## title

Milestone 17: model events run your scripts

## Links

### In

- derived-from [m17-picked-by-orchestrator](../decision/m17-picked-by-orchestrator.md)
- derived-from [event-hooks-from-daemon](../decision/event-hooks-from-daemon.md)

## History

- `01M6BV1J7WC184NBNWQCX7C9WJ` Plan milestone 17: model events run your scripts, picked by the orchestrator under the owner's standing instruction
