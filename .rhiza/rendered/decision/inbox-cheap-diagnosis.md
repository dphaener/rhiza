# inbox-cheap-diagnosis

decision · `01M6KVAAXE88YMEHYEZEEEFNB1`

## options

```json
[
  "read and fold the change sets once per inbox call, from the index when there is one",
  "cache inbox items in a daemon-side file keyed by log stamp",
  "compare overlap touched-field sets instead of what-if folds"
]
```

## quote

"We are just killing the CPU on this machine. All of these builds are taking way too long. This needs a diagnosis. I think we need to pause and figure this out." (owner, 2026-11-04)

## statement

Profile first, then fix what the profile shows: the inbox's per-lookup log reads (changeStateOf and changeTitleOf decoded the whole log once per run), the daemon fold's re-fold of the log for every candidate's model, and the daemon's per-poll re-parse of run journals and growing Claude Code transcripts. Overlaps themselves were not the hot spot.

## wording

the owner's words; the design is the orchestrator's, as the owner's proxy

## Links

### In

- decided-by [inbox-cheap](../capability/rhiza/orchestration/inbox/inbox-cheap.md)

## History

- `01M6KV89S5MKWG2DHA762F8WTZ` Inbox and daemon notify are cheap: overlaps and inbox items are computed once per model change, not per call or poll
