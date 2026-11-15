# m38-warm-at-ab67d485

evidence · `01M7FT15A5MDSHQ5WGA091BJ13` · actor [claude-code-8439813b](../actor/claude-code-8439813b.md)

## head

ab67d485c75e4c7529f28513cf048ba13745de31

## result

0:51.21 wall, 123% CPU, exit 0

## statement

At ab67d485c (all five milestone-38 change sets merged), a warm make test on the 4-core VM took 51.2 s of wall time (31 user + 32 sys CPU-s), within m38-warm's 60 s (2026-11-15: 254 s before). Measured by the orchestrator: make test, then /usr/bin/time make test.

## status

verified

## Links

### Out

- relates-to [m38-warm](../intent/milestone-38/m38-warm.md) — Evidence that m38-warm is met at ab67d485c
- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
