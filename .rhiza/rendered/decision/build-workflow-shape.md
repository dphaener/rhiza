# build-workflow-shape

decision · `01M5WPTMRWPNYM2AMVSMDREY04`

## quote

everything else looks good.

## quote_where

milestone-3 conversation, topic 4: the build workflow

## statement

In the build workflow a change set's leaves are the tasks: a leaf is ready when the leaves it depends-on have Evidence; a builder writes code and tests to the leaf's acceptance and a fresh verifier re-runs it and writes the Evidence; a failed check goes back to a builder up to 3 times on stronger models, then the leaf fails, its dependents are skipped and it becomes a question; outside work is a discovered-from change set or a question; one commit per leaf, and the build does not land.

## wording

the orchestrator's proposal in the milestone-3 conversation, accepted by the owner

## Links

### Out

- derived-from [milestone-3-conversation](../source/milestone-3-conversation.md)

### In

- decided-by [milestone-3](../intent/anyone-builds/milestone-3.md)
- decided-by [leaf-dependencies](../capability/rhiza/model/leaf-dependencies.md)
- decided-by [leaf-depends-on-acyclic](../capability/rhiza/model/leaf-dependencies/leaf-depends-on-acyclic.md)
- decided-by [leaf-order](../capability/rhiza/model/leaf-dependencies/leaf-order.md)
- decided-by [leaf-depends-on-guide](../capability/rhiza/model/leaf-dependencies/leaf-depends-on-guide.md)

## History

- `01M5WPTHR3EBHDNNH1CFX7A0T0` Milestone 3 planning record: the conversation, the milestone-3 intent and its done points, and its decisions
