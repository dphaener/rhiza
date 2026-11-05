# mac-banner-plan

source · `01M6JBEQH0K54X88TA8DRXGQZ8`

## quotes

```json
[
  {
    "by": "owner",
    "text": "I have seen no notification on my mac. But I have heard chimes",
    "where": "2026-11-03, change set intent"
  },
  {
    "by": "owner",
    "text": "Keep it open, bug first",
    "where": "2026-11-03, change set intent"
  }
]
```

## raw_sha256

31882d3e692dd728fa12fc7c5b540ee931e7dfc0ba529b489b498a57fc8e3cae

## session

b5fe1d11-4e23-4674-ae5d-01a1a4b11ff3

## summary

Autonomous plan run with no owner present. Investigation found that since the headless-run change (#210) run drivers start with no REX_BLOCK, so ctx.Rex is nil and rhiza's own blocked-run notice (root-record OSC 7501 with kind, title and msg) is never written anywhere; the only status Rex sees for a blocked agent tab is Claude Code's native report, which chimes on the Mac but shows no banner. Rex's program-status docs leave sounds and banners to the terminal and document no banner rules. Two questions went to the arbiter, who decided: the agent tab's own wrapper (rhiza session run) writes rhiza's blocked notice to its own terminal and a working report when the block ends, the driver writes to its own block only when it has one, and no driver, daemon or watcher tab is opened; after the fix the Mac is checked for a banner, and if it still only chimes a report for Rex upstream is drafted in the repo for the owner to file, done once drafted and pointed at in the inbox.

## title

Plan of the missing Mac banner for VM agent questions (autonomous)

## Links

### In

- derived-from [headless-blocked-notice-in-agent-tab](../decision/headless-blocked-notice-in-agent-tab.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
