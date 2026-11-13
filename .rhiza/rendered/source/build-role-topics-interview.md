# build-role-topics-interview

source · `01M7973DMXFE7BE30JE4EGSWVR`

## quotes

```json
[]
```

## raw_sha256

a59f14c0df9ae752bf18c4b5c60915cbfb20946bd4516c986726e11624a3d6be

## session

ab0574e9-12c9-4b2a-af3f-01a1d26b9855

## summary

Autonomous plan run for change set 01M76XR5 (Build role topics), serving m35-afk-prompts. No owner present. The planner read the change set intent, prompt-gap.md sections 2.3, 2.4 and 5, internal/guide (guide.go, build.md and its tests), internal/api/guide_install.go, internal/api/run_leaves.go leafPrompt, the builtin build.toml and review.toml verifier prompts and Rhiza's .rhiza/workflows/docs-build.toml. It raised one question bundling four choices (how shared blocks are written once, how the role topics reach sessions and what happens to the build topic, how a retry reaches build-retry, what happens to docs-build.toml). The arbiter answered with Decision build-role-topics-shape: shared blocks are files under internal/guide/text/shared/ expanded at load from `<!-- shared: <block> -->` lines; builder, verifier and build-retry are new topics installed as the rhiza-builder, rhiza-verifier and rhiza-build-retry skills; build.toml prompts shrink to one line; the build topic and rhiza-build skill stay as a short index (not breaking); retry is routed by the builder topic's TASK section on RHIZA_ATTEMPT > 1, with verbatim findings left to change set 01M76XT9; docs-build.toml's prompts become the builtin lines (plus its doc addendum for the builder). The planner then wrote six leaves under build-guide: shared blocks, the builder, verifier and build-retry topics, install plus the build index, and the one-line prompts (build.toml, review.toml's verifier, docs-build.toml).

## title

Build role topics plan (autonomous)

## Links

### In

- derived-from [build-role-topics-shape](../decision/build-role-topics-shape.md)
- derived-from [guide-shared-blocks](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/guide-shared-blocks.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
