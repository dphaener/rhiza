# roadmap-app-umbrella-row-qa-1

evidence · `01M6A312ZZ6GVYH4J1V9DH41TY` · actor [claude-code-36a5beca](../actor/claude-code-36a5beca.md)

## result

PASS

## statement

Ran the scratch binary's rhiza tui app roadmap in a real Rex terminal pane over the same repo and drove it with send-key. Collapsed, it showed '+ umb (2 change sets)' and '+ umb2 (1 change set)' with no member rows. j reached the umb row, and enter turned it into '- umb' with U and V listed below, with their not-ready reasons. j moved onto U, and enter opened U's detail. esc went back, and r kept the cursor on U with the group still expanded. k then enter collapsed it again, and r kept the cursor on the umb row. a, p, s and o on a header row do nothing. tui app roadmap --once printed the same bytes as rhiza roadmap --json, and cmd/rhiza/tui.go is unchanged. go test ./internal/tui/roadmapapp passes.

## Links

### Out

- verifies [roadmap-app-umbrella-row](../capability/rhiza/roadmap-app-list/roadmap-app-umbrella-row.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
