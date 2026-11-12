# watchdog-classifies-failures-qa-1

evidence · `01M77K5G7GH401Y3XS0RZR8CQ3` · actor [claude-code-097fafbe](../actor/claude-code-097fafbe.md)

## result

PASS

## statement

QA at 3a745c744 with the built binary: failed autonomous runs (rhiza run start --foreground --autonomous, step printing an HTTP 502 error) then rhiza daemon --once in a scratch root. Classifier empty: logged 'watchdog: run … step a failure classified real, left to the inbox: classifier off' and daemon.state.json records judged for the run. A fake 'transient: github 502' classifier received the failed step's detail on stdin; the watchdog logged 'failed on a transient error (github 502): resumed with --retry (budget 1 of 2)' and filed a bug request. A classifier exiting 2 logged 'classified real, left to the inbox: classifier failed: exit status 2', and a second daemon --once pass did not ask or log again. No transientSignatures remain in the code (grep).

## Links

### Out

- verifies [watchdog-classifies-failures](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-classifies-failures.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.
