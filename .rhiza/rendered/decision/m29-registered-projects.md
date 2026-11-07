# m29-registered-projects

decision · `01M6RN37C81BAPQKZZKWSEEK3W`

## question

What is the root directory, and how do projects get into it?

## statement

Projects are registered explicitly in the root's config with rhiza project add (a path, or a git URL to clone under the root); a clone that is not registered is never built.

## wording

Registered under root. The assistant's option, picked by the owner.

## Links

### Out

- derived-from [milestone-29-interview](../source/milestone-29-interview.md)

### In

- decided-by [m29-root-registry](../intent/milestone-29/m29-root-registry.md)
- decided-by [root-init](../capability/rhiza/runtime/root-init.md)
- decided-by [root-show](../capability/rhiza/runtime/root-show.md)
- decided-by [project-add-path](../capability/rhiza/runtime/project-add-path.md)
- decided-by [project-add-url](../capability/rhiza/runtime/project-add-url.md)
- decided-by [vm-root-installed](../capability/rhiza/orchestration/daemon/vm-root-installed.md)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
