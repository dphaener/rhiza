# fix-tests-write-fake-plan

source · `01M6T26DAP8QKQB2EWR4PNN74S`

## quotes

```json
[]
```

## raw_sha256

24d6824def9376b468bf802962a02d96aac4a668889a39baf6882cd477fac99b

## session

df8b7e8a-c713-40c2-ae5d-01a1b4206f96

## summary

Autonomous plan run for the change set fixing tests that wrote fake usage-limit hits into the machine's real root. The planner found the leak: testCtx leaves Ctx.Wd empty and rootcfg.Find("") resolves the process cwd, a package dir under ~/code. It linked the change set to m29-still-healthy and raised one question with two parts: where an empty Wd stops resolving a root, and how a test run fails on writes outside its temp dir. The arbiter decided (empty-wd-finds-no-root): rootcfg.Find("") reports no root, and TestMain in internal/api and cmd/rhiza snapshots any enclosing root's .rhiza-root/ state dir and fails the run on files that appeared or changed. No owner words; nothing left open.

## title

Fix tests write fake: autonomous plan interview

## Links

### In

- derived-from [empty-wd-finds-no-root](../decision/empty-wd-finds-no-root.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
