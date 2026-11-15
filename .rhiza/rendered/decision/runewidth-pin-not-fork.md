# runewidth-pin-not-fork

decision · `01M7EQF5W66CPHY1X3H3RPZT8H`

## reasoning

The brief allows pinning a version or deferring; go-runewidth's init runs as soon as a package importing it is linked, so rhiza cannot defer it without restructuring its imports, and a fork or replace would be a new dependency to maintain, which the planner rules forbid. v0.0.28 through v0.0.31 all init lazily; v0.0.31 is the latest and its go.mod requires the same github.com/clipperhouse/uax29/v2 v2.2.0 as v0.0.27, so the bump adds and moves no dependency.

## statement

rhiza's start stops paying go-runewidth's eager init by moving github.com/mattn/go-runewidth (indirect, through charmbracelet/colorprofile and charmbracelet/x/ansi) from v0.0.27 to v0.0.31, the latest release, whose init paints only the first 0x300 LUT entries and builds the rest lazily; never by a fork, a replace directive or a vendored copy, and with no other module's version changed

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)
- derived-from [cheap-start-plan-review](../source/cheap-start-plan-review.md)

### In

- decided-by [cheap-process-start](../capability/rhiza/integrity/self-check/cheap-process-start.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).
