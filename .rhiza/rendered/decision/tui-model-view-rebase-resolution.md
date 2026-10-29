# tui-model-view-rebase-resolution

decision · `01M63ZSVPQTN9X35WSZAH84WM2`

## reasoning

The two change sets each changed tuiApp's first lines: main split the lookup into read and write apps, this branch added a branch for the model app. The model app is neither a read nor a write app, so its check goes first and main's lookups stay as they are.

## statement

The rebase onto main is resolved in cmd/rhiza/tui.go by keeping main's read and write app lookups (from TUI write apps, previewed) and putting the model app's check before them: tuiApp hands the model app to tuiModel, then looks the name up among the read and write apps. Both changes keep their behaviour; rendered conflicts were settled by re-rendering, and make check passes on the rebased branch.

## wording

decided by the orchestrator acting as arbiter, under the owner's standing word to run autonomously

## Links

### Out

- answers [rebase-conflict-tui-model-view](../question/rebase-conflict-tui-model-view.md)

## History

- `01M62JMT63RY271EAF14R8P4VG` A Rex app showing the as-is capability tree and the roadmap
