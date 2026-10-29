# model-tree-app-qa-2

evidence · `01M6368ZM1HH1GP2MTKVFJXB5T` · actor [claude-code-e8b0da0f](../actor/claude-code-e8b0da0f.md)

## result

PASS

## statement

QA cycle 2 at 0184e67 with a scratch build of this branch, on gen.Default() (50,000-node model) written to a scratch git repo. `rhiza tui app model --once` equals `rhiza tree --view as-is` + blank line + `rhiza roadmap` byte for byte (cmp, 15,976 lines); with a deep-subtree ref (intent-57, an ancestor of the depth-6 sample) `rhiza tui app model intent-57 --once` equals `rhiza tree intent-57 --view as-is` + roadmap (118 lines). Whole-process wall time for the ref read: 0.03 s on 5 of 5 runs (no ref: 0.18 s, outside the budget per model-app-read-budget). Bad ref exits 4 with next: rhiza ls; extra args exit 2 with the usage line. No socket in the scratch repo. Interactive in a detached 80x24 Rex pane (`rhiza tui app model intent-30`): the screen is 24 lines with the [tree] roadmap header on line 1 and the key footer on the last line, the cycle-1 clipping fixed; j/j scroll with the header kept, tab switches to [roadmap], G jumps to the roadmap's last line, g back to # Roadmap, the header stays on screen throughout, q exits 0. go test ./internal/tui/modelapp and the tui-model-app script pass.

## Links

### Out

- verifies [model-tree-app](../capability/rhiza/model-tree-app.md)

## History

- `01M62JMT63RY271EAF14R8P4VG` A Rex app showing the as-is capability tree and the roadmap
