# model-tree-app-qa-1

evidence · `01M63581S7J20KZMHTN2ZG2FZR` · actor [claude-code-7cd94038](../actor/claude-code-7cd94038.md)

## result

PASS

## statement

QA at bcec030 with a scratch build of this branch: wrote gen.Default() (50,000-node model) to a scratch git repo. `rhiza tui app model --once` equals `rhiza tree --view as-is` + blank line + `rhiza roadmap` byte for byte (cmp; 15,938 tree lines, 37 roadmap lines); with the deep-subtree ref (grandparent of the deep sample) `rhiza tui app model <ref> --once` equals `rhiza tree <ref> --view as-is` + roadmap likewise (329 tree lines). Whole-process wall time for the ref read: 0.03 s on all 5 runs (whole tree, no ref: 0.17 s, the same as rhiza tree, which D14 benchmarks but does not budget). Same byte match on this repo's own .rhiza. Bad ref exits 4 with next: rhiza ls; extra args exit 2 with the usage line. No socket and no daemon lock in the scratch repo. Interactive in a detached 80x24 Rex pane: j/k scroll, tab switches to the roadmap, G/g jump to end and top, q exits 0. The [tree]/[roadmap] header line is clipped off the top (the view is 25 lines in 24 rows), reported as a finding.

## Links

### Out

- verifies [model-tree-app](../capability/rhiza/model-tree-app.md)

## History

- `01M62JMT63RY271EAF14R8P4VG` A Rex app showing the as-is capability tree and the roadmap
