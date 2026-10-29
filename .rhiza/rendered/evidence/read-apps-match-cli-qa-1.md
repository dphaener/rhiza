# read-apps-match-cli-qa-1

evidence · `01M62YY0CFTSHXNT9G1SVGZ88Z` · actor [claude-code-655e5878](../actor/claude-code-655e5878.md)

## result

PASS

## statement

QA at 4a03cf7 with a built rhiza on this repo's real model (RHIZA_ROOT/RUN unset). rhiza <op> --json and rhiza tui app <op> --once are byte-identical (cmp) for roadmap (19626 bytes), ready (14007 bytes) and slice with RHIZA_CHANGE=01M62JMSW4FH8RAWSQXHC7ZPND (18187 bytes); show/why/slice on handles not in the as-is view fail the same way (not_found, exit 4) in both. In a real Rex session (patched only to attach, see the qa finding) the roadmap modal showed the same JSON with the j/k scroll footer. Gap: tui app drops the global --change flag, so rhiza tui app slice --change <cs> --once fails (exit 2) where rhiza slice --change <cs> --json succeeds; reported as a minor finding.

## Links

### Out

- verifies [read-apps-match-cli](../capability/rhiza/read-apps-match-cli.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
