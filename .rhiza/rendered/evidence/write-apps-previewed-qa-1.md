# write-apps-previewed-qa-1

evidence · `01M636ARAZA0GEHPDN7RPCPJEF` · actor [claude-code-f19998b2](../actor/claude-code-f19998b2.md)

## result

PASS

## statement

QA at e74bb78 with a binary built from the branch (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza), run env unset, in scratch repos. Previews: tui app propose, record-add, record-set, record-link and record-retire --once each printed dry_run true, and failing previews (land on a draft, exit 3) wrote nothing; a hash of every .rhiza file was identical before and after. change-show --once output cmp-equal to change show --json. Confirm: in two identical repos (same RHIZA_TEST_NOW/SEED) I ran propose, accept --reason ok, land, record-add, record-set, record-link and record-retire with --once --confirm in one and the matching rhiza command with --json (--record for records) in the other: every output was cmp-equal and diff -r of the two .rhiza dirs showed no difference. Interactive in a real Rex pane: tui app record-add asked for arguments, showed the What if preview with dry_run true; pressing n exited 0 and the log hash was unchanged; a second run confirmed with y showed Written and the entry, and rhiza show cache then found the node.

## Links

### Out

- verifies write-apps-previewed (retired)

## History

- `01M62JMTG32M03T1FHC0T1EPE2` Change-set and record writes from apps, each shown as a what-if preview and confirmed once
