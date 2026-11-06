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

Per area (a package, feature or command group), write the capabilities its code and tests
deliver under a parent capability:

```
rhiza apply <file> --record --json
```

Keep one record per area, never one for the whole repo.

- Statement: one user-visible behaviour per capability, in the present tense.
- Acceptance: the repo's existing tests that prove it, named by test file (and test name). Untested: a question asking
  for the test, never an invented acceptance.
- Link a capability to its tests as evidence.

## 4. Bring in load-bearing decisions

Record decisions the code depends on (ADRs, comments, history), each `decided-by` its
capability. Never invent reasons: raise a question.

Bring in only a decision that shapes a capability or constraint; leave the rest.

A document's words are not the owner's: such a decision gets no `quote` field, this
wording form, and a link to its document's Source:

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

The `<file>` (`rhiza guide sources`, `rhiza guide scrub`) holds a scrubbed summary and the
owner's words, verbatim; `--raw` hashes the document.
Link what came from it with `derived-from`.

Keep reference, design and operations documents that capabilities link to: they stay in git.

## 6. Write unbuilt work as drafts

Planned work the code lacks: add its document as a Source, then a draft change set per
piece, naming it:

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

Fix what `check` reports, commit `.rhiza/rendered/`, spot-check capabilities against the
code.

## Record friction

Record every step you could not follow, or had to work around, as a question naming the
step, so the import can improve. Never work around one silently:

```
rhiza add question <handle> --field 'statement=Import step 3 (record built capabilities): <what, and how>'
```

## Follow-up: anchors

Once every area lands, one follow-up change set anchors each capability to its code:

```
rhiza add anchor <handle> --field path=<file> [--field symbol=<name>]
rhiza link <anchor> implements <capability>
```

`rhiza check` then fails on any anchor that does not resolve.
