# write-apps-previewed-verified

evidence · `01M6363269SMV955V143P6JW3V` · actor [claude-code-6a50202e](../actor/claude-code-6a50202e.md)

## result

PASS

## statement

At 74a4128 I ran go test ./cmd/rhiza -run TestScripts/tui-write-apps -count=1: PASS. The script cats every .rhiza/log file before and after --once previews of propose, record-add, record-set, record-link and record-retire (each prints dry_run true), and after failed previews of accept and record-retire, and cmp shows the log byte-identical. It then confirms propose, accept, land, record-add, record-set, record-link and record-retire with --once --confirm in repo a and runs the matching rhiza command with --json (and --record for record writes) in an identical repo b: each output cmps equal, the two full logs cmp equal and both fold-ok. change-show --once matches rhiza change show --json. go test ./internal/tui/writeapps passed (cancel with n makes only the dry-run call; y runs it for real once). go test ./..., fmt-check and vet passed; make self-check passed with the run env unset (with it set it failed only on this verifier session's own unrendered actor). I read the diff: write apps register with the launcher, run the api op in-process as a dry run first, and write only on one confirmation.

## Links

### Out

- verifies write-apps-previewed (retired)

## History

- `01M62JMTG32M03T1FHC0T1EPE2` Change-set and record writes from apps, each shown as a what-if preview and confirmed once
