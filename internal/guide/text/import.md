# Import

Bring an existing repo with built code into Rhiza. Seven steps, in this order.

## 1. Init

```
rhiza init
```

Run it at the repo root. It creates `.rhiza/`. Commit that on its own first.

## 2. Survey

Read before you write: the repo's code, tests and documents (READMEs, ADRs, design notes,
plans, TODOs). Note what is built, decided and only planned. Write nothing yet.

## 3. Record built capabilities, area by area

Take one area at a time (a package, a feature, a command group). Write the capabilities
its code and tests already deliver, under a parent capability for the area:

```
rhiza apply <file> --record --json
```

One `--record` change set per area keeps each reviewable: one record per area, never one for
the whole repo. Describe what the code does today, not what it should do.

- Statement: one user-visible behaviour per capability, in the present tense.
- Acceptance: the repo's existing tests that prove it, named by test file (and test name where
  there is one). A capability no existing test proves gets a statement, and a question asking
  for the test; never an invented acceptance.
- Link a capability to its tests as evidence where they prove it.

## 4. Bring in load-bearing decisions

Record the decisions the code depends on (from ADRs, comments, history) as decisions, each
hung off what it decides with `decided-by`. Never invent reasons: where the why is unknown,
raise a question instead.

The load-bearing rule: bring in only a decision that shapes a capability or constraint;
leave the rest in the document.

A document's words are not the owner's. A decision with no owner words (an agent wrote the
document) gets no `quote` field, this wording form, and a link to the Source of its document:

```
rhiza add decision <handle> --field 'statement=<decided>' \
  --field 'wording=<as the document states it>, made autonomously by the agent and not confirmed by the owner'
rhiza link <handle> derived-from <source>
```

## 5. Retire planning documents into Sources

A plan, roadmap or design note that has been built or overtaken becomes a Source
(`rhiza guide sources`; scrub it first with `rhiza guide scrub`). Link the decisions and
capabilities that came from it with `derived-from`. Then delete or archive the file in git.

## 6. Write unbuilt work as draft change sets

What the documents plan but the code lacks is not recorded as built. Open a draft change set
per piece of work (`rhiza change new "<intent>" --json`), write its intent and capabilities
there, and leave it in `draft` for the owner to propose.

## 7. Verify

```
rhiza check
rhiza render
```

Fix what `check` reports, commit `.rhiza/rendered/`, and spot-check a few capabilities: the
code does what they say, and `rhiza why` shows their decisions.

## Follow-up: anchors

Record the capabilities first and anchor them after. Once every area is landed, open one
follow-up change set and, for each capability, add an anchor to the code that implements it
and link it:

```
rhiza add anchor <handle> --field path=<file> [--field symbol=<name>]
rhiza link <anchor> implements <capability>
```

Then `rhiza check` again: it fails when an anchored path or symbol does not resolve.
