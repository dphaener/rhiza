# plan-roadmap-row-titles

source · `01M79TE4N9H2VSZN4C6CJQQFVQ`

## quotes

```json
[
  {
    "by": "owner",
    "text": "The titles seem to be a description and just run off. It's not informative at all.",
    "where": "bug request, 2026-11-13"
  }
]
```

## raw_sha256

2b68fb3baac07b048190df0ab02c465b9b869485e451900fe9300ce269d7917f

## session

f9012d6a-ac8d-4507-910c-01a1d3a461d7

## summary

Autonomous plan of a bug request: the web roadmap row title is the change set's whole intent (internal/web/roadmap.go sets Row.Title to c.Intent; the roadmap operation's change sets carry no title). Planned one leaf under web-roadmap-page serving m37-roadmap (decision web-redesign-direction: clear roadmap bugs ship at once as bug requests): the roadmap operation adds each change set's display title, and the web roadmap row shows it with the intent as a secondary line clamped to two lines. The runs list already shows the change set's display title and the change set page's depends-on list already uses change.show's title, so they need no change; the overview lists milestones, not change sets. No question was raised.

## title

Plan: web roadmap rows show change set titles

## Links

### In

- derived-from [roadmap-rows-by-title](../decision/roadmap-rows-by-title.md)

## History

- `01M79T7NW1B0KKXVVSX4AHZNHF` Fix The web roadmap shows each change set's whole intent sentence as its title, so rows run off and say little (owner, 2026-11-13: 'The titles seem to be a description and just run off. It's not informative at all.'). roadmap.go:140 sets Row.Title to c.Intent; change sets carry a short title (rhiza change new --title, at most 30 characters, defaulting to the intent's first clause). Wanted: rows show the change set's title, with the intent as the secondary line clamped to two lines, on the roadmap and wherever else the web UI lists change sets (overview, runs list), with a test.
