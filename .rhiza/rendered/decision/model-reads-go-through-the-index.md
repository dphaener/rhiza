# model-reads-go-through-the-index

decision · `01M6NEHG641TR9VC8NXJ78M4GA`

## reasoning

The change set's intent asks to poll cheaply (file mtimes or the run journal) and cache the fold. internal/index already does both: it is a disposable SQLite index derived from the log, catch-up is input-incremental (a seen entry path is never read again), it stores each entry's canonical bytes and each change set's folded state, and Open has a stat fingerprint so the fresh case costs one read of log/ and a stat per folder. Read operations already use it through openReader. The ~5 s read and ~1.5 s fold the question measured come from the paths that still call workspace.log() / store.ReadLog (writes such as apply, change show, the ship check subprocess, driver polls). Routing those through the index fixes the cost where it is paid, reuses tested freshness and rebuild rules, and fits CLAUDE.md's direction that the interim FoldLog model is to be replaced by projection. A separate fold cache would duplicate the index and its invalidation rules. Full reads stay where correctness needs every byte: check (finds edited entries, which the index does not notice) and render.

## statement

Model reads in drivers, 'rhiza change show', 'rhiza ship check --settled' and 'rhiza apply' go through the existing per-worktree SQLite index (internal/index) instead of reading and folding the whole log. Catch-up reads only entry files it has not seen, and is made incremental where it is not yet: today refold decodes every stored entry and runs changeset.FoldLog over the whole log on any catch-up that reads something new, and the change_sets table holds no history; this change set makes catch-up refold only the change sets its new entries touch and store their folded history, so stored canonical bytes and folded change sets are reused rather than decoded and re-canonicalized again. A poll whose stat fingerprint is unchanged does no model work, and a poll after another run appends an unrelated entry costs only that entry's change set. Writes validate against the index's folded state plus their own new entries. Of the commands this change set covers (driver polls, change show, ship check --settled, apply and the single-write commands), none reads the whole log; 'rhiza check' and 'rhiza render' keep reading it by design, and other reads that still call workspace.log() (status, inbox, notes show, change diff, carve, the daemon fold) are left for later change sets. No second cache (no separate fold file or in-memory fold kept across commands) is added. An idle driver gets a CPU budget test

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-a750dc81-2](../question/question-a750dc81-2.md)
- derived-from [cpu-drivers-plan-interview](../source/cpu-drivers-plan-interview.md)

### In

- decided-by [run-driver-idle-poll-cheap](../capability/rhiza/orchestration/daemon/run-driver-idle-poll-cheap.md)
- decided-by [change-show-ship-check-from-index](../capability/rhiza/views/fast-reads/change-show-ship-check-from-index.md)
- decided-by [apply-validates-from-index](../capability/rhiza/views/fast-reads/apply-validates-from-index.md)
- decided-by [index-catchup-refolds-incrementally](../capability/rhiza/views/fast-reads/index-catchup-refolds-incrementally.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
