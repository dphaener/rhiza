# init-asks-agent-setting

decision · `01M7E5H441QVC30ZVTC4A5V573`

## question

Change set 01M79KHH also says 'rhiza init asks which agent the project uses'. Today init is non-interactive, there is no project agent setting, and run start / session start / skill check take --agent with the built-in default adapter. How should init ask and what does the answer do?

## reasoning

Asking without storing the answer would change nothing a run does; a project setting in git config matches where the other project settings (rhiza.ship.*, rhiza.autonomy.*) live and is read the same way by run start, session start and skill check, which today fall back to adapter.Default(). Taking the default silently without a terminal keeps scripted inits and the testscripts working unchanged. Listing adapter.Names() keeps the agents-named-only-in-adapters constraint.

## statement

rhiza init takes --agent <adapter>; with no flag and a terminal on stdin it lists the registered adapters and asks, the default preselected (an empty answer takes it); with no terminal (scripts, --json) it takes the default without asking. The answer is stored in a new project setting rhiza.project.agent (git config), which run start, session start and skill check use when --agent is unset; an explicit --agent still wins. Adapter names come from the adapter registry, so nothing outside an adapter names an agent.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [byo-docs-planning](../source/byo-docs-planning.md)

### In

- decided-by [project-agent-setting](../capability/rhiza/runtime/agent-adapters/project-agent-setting.md)
- decided-by [init-asks-agent](../capability/rhiza/model/repo-init/init-asks-agent.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
