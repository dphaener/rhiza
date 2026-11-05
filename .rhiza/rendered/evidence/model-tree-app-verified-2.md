# model-tree-app-verified-2

evidence · `01M63647HG1GD3A5GJ8WYSC12A` · actor [claude-code-d5626906](../actor/claude-code-d5626906.md)

## result

PASS

## statement

Verifier of attempt 2 at 755b161: wrote gen.Default() (the 50,000-node model) to a scratch repo and ran rhiza built from this commit. `rhiza tui app model --once` equals `rhiza tree --view as-is` + blank line + `rhiza roadmap` byte for byte (cmp; 15,938 tree lines, 37 roadmap lines), and `rhiza tui app model <deep-subtree-ref> --once` equals `rhiza tree <ref> --view as-is` + roadmap likewise. The subtree+roadmap read, the budget model-app-read-budget defines, takes 0.03 s whole-process over 5 runs, and TestGeneratedModel's in-process median is 8.9-12 ms (< 50 ms); with no ref the whole tree takes 0.16 s, outside that budget as recorded. The fix for qa f1 holds: rows = Height-4, TestViewFitsTheWindow passes, and in a live 24-row Rex pane the capture is 24 lines with the [tree]/[roadmap] tab bar on line 1 after tab and j scrolling. fmt-check, vet and go test ./... pass (one full-suite run first failed TestGeneratedModel at a 78 ms median under parallel load; it passed on rerun and alone). make self-check fails only on stale .rhiza/rendered files, which were already stale at 080586f before this leaf's fix commit and are left for the orchestrator's render. Diff read: a one-line row-count fix, a test, and the model-app-read-budget Decision.

## Links

### Out

- verifies model-tree-app (retired)

## History

- `01M62JMT63RY271EAF14R8P4VG` A Rex app showing the as-is capability tree and the roadmap
