# milestone-25

intent · `01M6HKVRE1B7RPP12EN6J2GGXC`

## milestone

25

## statement

The backlog is empty: every change set planned from the owner's queued requests and the backlog is built and shipped by the daemon on its own, so the owner can start on new things, bugs and polish

## Links

### Out

- decided-by [m25-empty-backlog](../decision/m25-empty-backlog.md)

### In

- serves change set `01M6HKTEKH6WQTCXRYDM9GSASG` (Milestone 25: empty the backlog — every planned change set from tonight's requests serves it so the daemon builds them)
- serves change set `01M6HHMTEDFK2AWQKX50EB77DD` (Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate)
- serves change set `01M6H3A65RDP761HKE82STZSQZ` (log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan))
- serves change set `01M6HJXAJ6GYA2RTAGVK6JR1N9` (Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them)
- serves change set `01M6HKRVC93TGX0ZC40QZ2GW9J` (Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there)
- serves change set `01M6HR3VBT4DPCNS03S70QCAZF` (Timing budgets measure in quiet: make check runs TestReadBudget, TestBudget and TestGeneratedModel on their own after the rest of the suite, one package at a time, so CI and the VM stop failing them on their neighbours' load)
- serves change set `01M6HN24SPSS76PKZHZV188JEW` (Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows))
- serves change set `01M6HNNWFR1J2AHQRNKFTPME46` (Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model))
- serves change set `01M6HSZMD2H2HZVGQQSVZZ4A1D` (Plan requests are not starved by ready change sets: while a plan request waits, the daemon keeps one run slot for it)
- serves change set `01M6JJMK3M50P34PPX1QNTHYCB` (The test suite is cheap to run in many worktrees at once: unchanged packages hit the shared test cache across worktrees, and concurrent full suites on one machine queue instead of thrashing)
- serves change set `01M6K044DTSK7PT604NHPNGD17` (Queued continuations are not starved by ready change sets: the daemon keeps a slot for each queued resume or retry)
- serves change set `01M6K3QFYNKE35DZJG5HPHPJT3` (Order the inbox fix after the resume-after-land fix: both rewrite the same resume and inbox nodes)

## History

- `01M6HKTEKH6WQTCXRYDM9GSASG` Milestone 25: empty the backlog — every planned change set from tonight's requests serves it so the daemon builds them
