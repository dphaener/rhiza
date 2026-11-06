# m29-one-root-daemon

decision · `01M6RN37AH4PP5N792VVNZF0S5`

## question

How should one machine run several projects? Today each repo has its own daemon, web UI port, usage governor and run cap.

## statement

One rhiza daemon per machine runs on the root and drives every project's runs under one usage governor and one machine-wide run cap, serving one web UI for all projects.

## wording

One root daemon. The assistant's option, picked by the owner.

## Links

### Out

- derived-from [milestone-29-interview](../source/milestone-29-interview.md)

### In

- decided-by [milestone-29](../intent/milestone-29.md)
- decided-by [m29-one-daemon](../intent/milestone-29/m29-one-daemon.md)
- decided-by [root-daemon-lock](../capability/rhiza/orchestration/daemon/root-daemon-lock.md)
- decided-by [root-daemon-polls-projects](../capability/rhiza/orchestration/daemon/root-daemon-polls-projects.md)
- decided-by [root-daemon-usage-governor](../capability/rhiza/orchestration/usage-governor/root-daemon-usage-governor.md)
- decided-by [root-daemon-service](../capability/rhiza/orchestration/daemon/root-daemon-service.md)
- decided-by [root-daemon-self-update](../capability/rhiza/orchestration/daemon/root-daemon-self-update.md)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
