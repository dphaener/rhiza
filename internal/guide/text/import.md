# Import

Bring an existing repo into Rhiza. Seven steps, in order.

## 1. Init

```
rhiza init
```

Run it at the root; commit `.rhiza/` first.

## 2. Survey

Read the code, tests and documents (READMEs, ADRs, plans, TODOs). Note what is built,
decided and planned. Write nothing yet.

## 3. Record built capabilities, area by area

Take one area at a time (a package, feature or command group). Write the capabilities its
code and tests deliver, under a parent capability:

```
rhiza apply <file> --record --json
```

One `--record` change set per area keeps each reviewable: one record per area, never one for
the whole repo.

- Statement: one user-visible behaviour per capability, in the present tense.
- Acceptance: the repo's existing tests that prove it, named by test file (and test name). A capability no test proves gets a statement and a question asking for
  the test; never an invented acceptance.
- Link a capability to its tests as evidence.

## 4. Bring in load-bearing decisions

Record the decisions the code depends on (ADRs, comments, history), each `decided-by`
its capability. Never invent reasons: raise a question.

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

## 6. Write unbuilt work as drafts

Planned work the code lacks is not built. Add its document as a Source, then open a draft
change set per piece, naming it:

```
rhiza change new "<intent>" --source <source> --json
```

Leave it in `draft`: accepting it is the owner's, never the import's. Write each open
question (a "TBD") as a question node, never a guess:

```
rhiza add question <handle> --field 'statement=<question>'
```

## 7. Verify

```
rhiza check
rhiza render
```

Fix what `check` reports, commit `.rhiza/rendered/`, spot-check a few capabilities
against the code.

## Record friction

Throughout the import, record every step you could not follow, or had to work around, as a
question in the importing repo's model, naming the step, so the import can be improved from
what it found. Never work around a step silently, and never stop the import over one:

```
rhiza add question <handle> --field 'statement=Import step 3 (record built capabilities): <what could not be followed or was worked around, and how>'
```

## Follow-up: anchors

Once every area is landed, open one follow-up change set and, per capability, add an
anchor to its code and link it:

```
rhiza add anchor <handle> --field path=<file> [--field symbol=<name>]
rhiza link <anchor> implements <capability>
```

Then `rhiza check`: it fails when an anchored path or symbol does not resolve.
