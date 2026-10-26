# agent-step-has-prompt

decision · `01M5WFCA132QY27VP37RDEKXD6`

## quote

an agent step has a `prompt` in the workflow file: a text template, passed to the adapter as the agent's first prompt.

## quote_where

the orchestrator's fix brief for change set 4 (m2-cs4-fixes.md), item 1

## statement

An agent step has a prompt in the workflow file, a text template over {change}, {slug}, {title}, {run} and {step}, passed to the adapter as the agent's first prompt; a step without one, or with an unknown placeholder, is invalid, and skills are invoked through the prompt text.

## wording

the orchestrator's ruling, made under the owner's authorization to amend plans (orchestrator-amends-plans); not the owner's words

## Links

### Out

- answers [agent-step-prompt](../question/agent-step-prompt.md)

### In

- decided-by [agent-step](../capability/rhiza/workflows/workflow-definitions/agent-step.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
