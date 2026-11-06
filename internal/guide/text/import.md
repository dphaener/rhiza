# Import

Bring an existing repo with built code into Rhiza. Seven steps, in this order.

## 1. Init

```
rhiza init
```

Run it at the repo root; commit `.rhiza/` on its own first.

## 2. Survey

Read the code, tests and documents (READMEs, ADRs, plans, TODOs). Note what is built,
decided and only planned. Write nothing yet.

## 3. Record built capabilities, area by area

Take one area at a time (a package, feature or command group). Write the capabilities its
code and tests already deliver, under a parent capability for the area:

```
rhiza apply <file> --record --json
```

One `--record` change set per area keeps each reviewable: one record per area, never one for
the whole repo. Describe what the code does today.

- Statement: one user-visible behaviour per capability, in the present tense.
- Acceptance: the repo's existing tests that prove it, named by test file (and test name where
  there is one). A capability no test proves gets a statement and a question asking for
  the test; never an invented acceptance.
- Link a capability to its tests as evidence.

## 4. Bring in load-bearing decisions

Record the decisions the code depends on (ADRs, comments, history), each hung off what it
decides with `decided-by`. Never invent reasons: where the why is unknown, raise a question.

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

Retire a built or overtaken planning document (PLAN, roadmap, TODO, design note): add it as
a Source and delete it in the same change set.

```
rhiza source add <handle> <file> --raw <document> --json
git rm <document>
```

The `<file>` (`rhiza guide sources`, scrubbed per `rhiza guide scrub`) holds a scrubbed
summary and the owner's typed words in the document, verbatim; `--raw` hashes the document.
Link what came from it with `derived-from`.

Keep reference, design and operations documents that capabilities link to: they stay in git.

## 6. Write unbuilt work as draft change sets

Planned work the code lacks is not recorded as built. Open a draft change set per piece
(`rhiza change new "<intent>" --json`), write its intent and capabilities there, and leave
it in `draft` for the owner to propose.

## 7. Verify

```
rhiza check
rhiza render
```

Fix what `check` reports, commit `.rhiza/rendered/`, and spot-check a few capabilities
against the code and `rhiza why`.

## Follow-up: anchors

Once every area is landed, open one follow-up change set and, per capability, add an
anchor to its code and link it:

```
rhiza add anchor <handle> --field path=<file> [--field symbol=<name>]
rhiza link <anchor> implements <capability>
```

Then `rhiza check`: it fails when an anchored path or symbol does not resolve.
