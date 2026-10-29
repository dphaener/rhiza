# rhiza-tui

intent · `01M62JMHNS9HS0AR6DT5SZYKGQ`

## statement

Rhiza's commands run as small Rex apps: each app is its own Rex action in the command palette and a key in the rhiza key mode, and opens straight away in a layer sized to it with no picker in between; each app reads and writes the model in-process through internal/api

## Links

### Out

- decided-by [tui-many-mini-apps](../decision/tui-many-mini-apps.md)

### In

- serves change set `01M62JMSW4FH8RAWSQXHC7ZPND` (Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps)
- serves change set `01M62JMT63RY271EAF14R8P4VG` (A Rex app showing the as-is capability tree and the roadmap)
- serves change set `01M62JMTG32M03T1FHC0T1EPE2` (Change-set and record writes from apps, each shown as a what-if preview and confirmed once)
- serves change set `01M62JMTT1Z8FP0SC619V3QS68` (Read-only Rex apps for run status and session ls)

## History

- `01M62J0ED188JTTR7R9709Y989` Rhiza has a terminal UI built on Rex, with rexuikit, in which common rhiza commands run as tiny apps in modals or pop-ups, alongside a view of the model.
- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
