# build-ruling-3

decision · `01M5ZARSENEQS7RK867GMCN97Y`

## options

```json
[
  "Take the old pane's place",
  "New tab each builder"
]
```

## picked

Take the old pane's place

## question

Where does the next builder's pane go once the previous builder passes and its pane closes?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

Where does the next builder's pane go once the previous builder passes and its pane closes? — Take the old pane's place. The next builder takes the passed builder's pane: the engine splits the old builder pane, then closes the old one, and opens a new tab only when the change set's tab is gone.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [build-panes](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-panes.md)
- supersedes [run-window-gone-opens-window](run-window-gone-opens-window.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
