# scoped-tests-plan

source · `01M6M0DB1S00W14FAP15T97J8T`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge.",
    "where": "2026-11-04, the change set's opening request"
  }
]
```

## raw_sha256

ee12b52948a7cd532a4dc0571299c03d4f38ddf42d7acd244669a1c9af50a972

## session

c2063c96-e994-4af5-babc-01a1a800e15c

## summary

Autonomous plan of change set 01M6M01G. The owner asked, when opening it, that agents in runs (builders, verifiers, review checks, critics, fixers) run only the tests their change touches, since the PR's GitHub CI runs the full suite before every merge and the VM's CPU is over-committed (~35% steal). The planner found that only this repo's CLAUDE.md tells agents to run make check; the build, review and plan guide topics (the rhiza-build, rhiza-review and rhiza-plan skills) and the built-in workflows' role prompts say nothing about test scope. Two questions went to the arbiter: which intent the change set serves, and how specific the guides' wording should be. The arbiter ruled (Decision scoped-tests-serve-m25-generic-rule): it serves milestone-25, and the guides and role prompts state the rule repo-agnostically, naming no Go, make or testscript command; this repo's concrete commands go in CLAUDE.md and the leaves' acceptance checks. Four leaves: the build/review guides, the plan guides' acceptance advice, the built-in role prompts, and CLAUDE.md. Nothing is left open.

## title

Plan: agents in runs stop running the full test suite

## Links

### In

- derived-from [ci-runs-full-suite](../decision/ci-runs-full-suite.md)
- derived-from [scoped-tests-serve-m25-generic-rule](../decision/scoped-tests-serve-m25-generic-rule.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
