# model-app-read-budget

decision · `01M635BW4PX2TGH5J2GM4S089T`

## statement

The model app takes an optional [<ref>] (`rhiza tui app model [<ref>]`) that limits the tree to the subtree under that node, as `rhiza tree <ref>` does. model-tree-app's 50 ms budget covers the app's read of a deep subtree plus the roadmap, the tree read that t7-d14-read-budget budgets; with no ref the app reads the whole as-is tree, which is outside the budget, as a whole-tree `rhiza tree` is.

## wording

fixer of review cycle 2 on leaf model-tree-app, recording the reading of the budget the build and its evidence used (model f1)

## Links

### Out

- relates-to [t7-d14-read-budget](t7-d14-read-budget.md) — applies t7-d14-read-budget's 50 ms tree-read budget to the model app's subtree read

### In

- decided-by [model-tree-app](../capability/rhiza/model-tree-app.md)

## History

- `01M62JMT63RY271EAF14R8P4VG` A Rex app showing the as-is capability tree and the roadmap
