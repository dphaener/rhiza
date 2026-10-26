# milestone-1

source · `01M5SWRGKRMVD1X8MV4CS5HB3Q`

## quotes

```json
[]
```

## raw_sha256

5c45d80082b961155e86c760cd4e821430ac6db5f2059f5c8b4615262b721914

## summary

The temporary task list for Genesis milestone 1, the only plan document Genesis was to have; task 11 brings it in as a Source and retires it, and the reasons behind it are in the founding conversation. Goal: the smallest Genesis that can take in the founding conversation and then plan its own next milestone from inside itself, so no PLAN.md ever exists in its repo. Go, one binary; built with Claude Code in Rex and shipit, one task per run and per PR, supervised on the Mac and never with permissions bypassed. In: the store (init, repo.json, IDs and handles, the log of immutable entries with the envelope, six operations, canonical JSON, append-only; change sets with revisions and the whole lifecycle; projection with as-is, to-be, what-if, causal order and conflict detection; the SQLite index per worktree, rebuilt incrementally); the writes (change new/use/show/ls, add, set, mv, retire, link, unlink, apply, --record, propose, accept, reject, withdraw, land); the reads (show, tree, ls, slice, why, log, roadmap, ready); integrity (render with committed markdown, check, ci install); agents (operations with schemas, describe, api call, --json, exit codes, next in errors, the guide, the genesis and genesis-scrub skills, guide install for Claude Code with the session-start hook, the size-cap test, source add); tests (testscript, golden files, property tests, the format-1 corpus, the 50,000-node generator and benchmark). Out, for milestone 2 or later: anchors, fingerprints, trailers and drift; reconcile; edit through $EDITOR; rebase; full-text find; repos and cross-repo links; hooks; the change feed; notes; the personal layer; AGENTS.md install; evals; every layer from runtime up. Done means: (1) the import: one landed record change set and the twelve areas as accepted intention change sets in the agreed order; (2) roadmap shows the founding plan, why gives back the owner's words and source, ready names the next change set, and its change slice fits 4,000 tokens with the root constraints; (3) make check, including genesis check on its own repo, passes in CI; (4) reads under 50 ms on the 50,000-node model; (5) milestone 2 planned and carried out inside Genesis. Tasks and merges: 1 the repo, Go module, make check, CI, testscript harness, CLAUDE.md (PR #1, 5a4b716); 2 IDs, canonical JSON, envelope, atomic writes, log reader, corpus (#2, 396505b); 3 operations, folding revisions, the lifecycle (#3, aea0f85); 4 projection and property tests (#4, 380f2f2); 5 the SQLite index, generator and benchmark (#7, ed8d292); 6 operations and the CLI frame plus every write (#5, 6d5c61e); 7 the reads (#10, 29e811d; #11, 3b5d1da); 8 slice (#9, a951e38); 9 render, check, ci install, self-check (#6, 8355823); 10 the guide, both skills, guide install, size cap, source add (#8, 3e17d7b); 11 the import, owner and agent, in progress when imported.

## title

Milestone 1 task list (temporary, retired by the import)

## Links

### In

- derived-from [node-ids](../capability/genesis/model/node-ids.md)
- derived-from [operation-log](../capability/genesis/model/operation-log.md)
- derived-from [format-compat](../capability/genesis/model/format-compat.md)
- derived-from [node-kinds](../capability/genesis/model/node-kinds.md)
- derived-from [typed-links](../capability/genesis/model/typed-links.md)
- derived-from [handles](../capability/genesis/model/handles.md)
- derived-from [six-operations](../capability/genesis/model/six-operations.md)
- derived-from [change-lifecycle](../capability/genesis/change-sets/change-lifecycle.md)
- derived-from [revisions](../capability/genesis/change-sets/revisions.md)
- derived-from [capability-tree](../capability/genesis/model/capability-tree.md)
- derived-from [conflicts](../capability/genesis/change-sets/conflicts.md)
- derived-from [overlays](../capability/genesis/views/overlays.md)
- derived-from [fast-reads](../capability/genesis/views/fast-reads.md)
- derived-from [repo-init](../capability/genesis/model/repo-init.md)
- derived-from [current-change](../capability/genesis/change-sets/current-change.md)
- derived-from [record-shortcut](../capability/genesis/change-sets/record-shortcut.md)
- derived-from [model-writes](../capability/genesis/change-sets/model-writes.md)
- derived-from [apply-files](../capability/genesis/change-sets/apply-files.md)
- derived-from [actors](../capability/genesis/lineage/actors.md)
- derived-from [command-schemas](../capability/genesis/guide/command-schemas.md)
- derived-from [why-walk](../capability/genesis/lineage/why-walk.md)
- derived-from [change-history](../capability/genesis/lineage/change-history.md)
- derived-from [node-show](../capability/genesis/views/node-show.md)
- derived-from [tree-view](../capability/genesis/views/tree-view.md)
- derived-from [node-list](../capability/genesis/views/node-list.md)
- derived-from [roadmap](../capability/genesis/views/roadmap.md)
- derived-from [ready-next](../capability/genesis/views/ready-next.md)
- derived-from [paged-lists](../capability/genesis/views/paged-lists.md)
- derived-from [slice](../capability/genesis/views/slice.md)
- derived-from [check](../capability/genesis/integrity/check.md)
- derived-from [render](../capability/genesis/integrity/render.md)
- derived-from [ci-install](../capability/genesis/integrity/ci-install.md)
- derived-from [self-check](../capability/genesis/integrity/self-check.md)
- derived-from [source-add](../capability/genesis/lineage/source-add.md)
- derived-from [scrubbing](../capability/genesis/lineage/scrubbing.md)
- derived-from [agent-guide](../capability/genesis/guide/agent-guide.md)
- derived-from [skills](../capability/genesis/guide/skills.md)
- derived-from [guide-install](../capability/genesis/guide/guide-install.md)
- derived-from [session-start](../capability/genesis/guide/session-start.md)
- derived-from [guide-size-cap](../capability/genesis/guide/guide-size-cap.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
