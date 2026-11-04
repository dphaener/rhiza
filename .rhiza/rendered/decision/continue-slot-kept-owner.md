# continue-slot-kept-owner

decision · `01M6K05ZW2H5WPK3SYCBMQGV3E`

## statement

Ship picks leave a slot for each queued continuation that could start now, plus the one kept for a waiting plan request: two runs whose open-pr failed on a GitHub 500 were queued for retry and never started while 9 change sets were ready, and raising max-runs gave the new slot to another ship pick. Continuations are short (open-pr, a resume after a wait), so the kept slots cost little

## wording

the orchestrator's, as the owner's proxy (owner, 2026-11-04: "i think you should monitor the daemon and adjust accordingly. we likely have some performance issues to tackle.")

## Links

### Out

- relates-to [request-slot-kept-owner](request-slot-kept-owner.md) — Generalises the plan request's kept slot to queued continuations

### In

- decided-by [continue-slot-kept](../capability/rhiza/orchestration/daemon/daemon-picks/continue-slot-kept.md)

## History

- `01M6K044DTSK7PT604NHPNGD17` Queued continuations are not starved by ready change sets: the daemon keeps a slot for each queued resume or retry
