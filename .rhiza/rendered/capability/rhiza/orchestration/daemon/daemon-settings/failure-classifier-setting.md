# failure-classifier-setting

capability · `01M7642N1KXSN7CKCDEZFBK78Q` · parent [daemon-settings](../daemon-settings.md)

## acceptance

Unit tests run a fake classifier command for transient, real, garbage output and a failing command, and show the parsed verdict, with real assumed on anything but a clear transient; the daemon-settings testscript lists the setting with its default; a testscript whose command step fails (e.g. build-leaf-fails) runs with the classifier empty from Setup and records no classifier call, and a script that sets a fake classifier command sees its verdict used.

## statement

A machine-wide setting names the failure classifier command: it reads a failed step's output on stdin and prints whether the failure is transient or real with a one-line reason; the default runs Claude headless with Haiku, any LLM CLI can be configured, and empty turns classification off. The testscript Setup in cmd/rhiza/main_test.go sets the classifier empty, so no test calls an agent or the network; scripts that test classification set a fake classifier command.

## Links

### Out

- decided-by [d-llm-classifier](../../../../../decision/d-llm-classifier.md)

### In

- verifies [failure-classifier-setting-verified](../../../../../evidence/failure-classifier-setting-verified.md)
- verifies [failure-classifier-setting-qa-1](../../../../../evidence/failure-classifier-setting-qa-1.md)
- verifies [failure-classifier-setting-verified-2](../../../../../evidence/failure-classifier-setting-verified-2.md)
- verifies [failure-classifier-setting-qa-2](../../../../../evidence/failure-classifier-setting-qa-2.md)
- verifies [pr-480-failures-classified-by-a](../../../../../evidence/pr-480-failures-classified-by-a.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.
