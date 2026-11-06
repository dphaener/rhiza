# Import

Bring an existing repo with built code into Rhiza. Seven steps, in this order.

## 1. Init

```
rhiza init
```

Run it at the repo root. It creates `.rhiza/`. Commit that on its own before anything else.

## 2. Survey

Read before you write. List the repo's code (packages, entry points), its tests and its
documents (READMEs, ADRs, design notes, plans, TODOs). Note what is built, what is decided
and what is only planned. Write nothing yet.

## 3. Record built capabilities, area by area

Take one area at a time (a package, a feature, a command group). For each, write the
capabilities its code and tests already deliver, under a parent capability for the area:

```
rhiza apply <file> --record --json
```

One `--record` change set per area keeps each reviewable. Describe what the code does today,
not what it should do. Link a capability to its tests as evidence where they prove it.

## 4. Bring in load-bearing decisions

Record the decisions the code depends on (from ADRs, comments, history) as decisions, each
hung off what it decides with `decided-by`. Quote the document's own words. Never invent
reasons: where the why is unknown, raise a question instead.

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

Fix what `check` reports, commit `.rhiza/rendered/`, and spot-check: pick a few capabilities
and confirm the code does what they say, and that `rhiza why` shows their decisions.
