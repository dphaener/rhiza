# milestone-2-conversation

source · `01M5SZHJ8S6TDVVJGKSMZGC3SG`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Yes let's do that",
    "where": "milestone-2 conversation, start: answering \"Shall we start planning milestone 2?\""
  },
  {
    "by": "owner",
    "text": "A",
    "where": "milestone-2 conversation, topic 1: what milestone 2 is for"
  },
  {
    "by": "owner",
    "text": "Yes to all",
    "where": "milestone-2 conversation, topic 2: how milestone 2 lives in the model, and the queue"
  },
  {
    "by": "owner",
    "text": "Yes to all 5",
    "where": "milestone-2 conversation, topic 3: change sets that build on each other"
  },
  {
    "by": "owner",
    "text": "Yes to all",
    "where": "milestone-2 conversation, topic 4: what Genesis runs the work covers"
  },
  {
    "by": "owner",
    "text": "Yes to all",
    "where": "milestone-2 conversation, topic 5: milestone 2 as change sets"
  },
  {
    "by": "owner",
    "text": "Let's go with a, and yes you have my approval.",
    "where": "milestone-2 conversation, after the rehearsal: proposing the drafts, and the go-ahead"
  },
  {
    "by": "owner",
    "text": "Ok let's add the recommended change set",
    "where": "milestone-2 conversation, after carving was found not to work as decided"
  },
  {
    "by": "owner",
    "text": "All of them, let's go",
    "where": "milestone-2 conversation, topic 6: carving that keeps IDs"
  },
  {
    "by": "owner",
    "text": "Ok I think I like the name. The designs are going to need some work but we've got a good start. I'm going to table that for now, but let's move forward with the name change.",
    "where": "milestone-2 conversation, the name"
  },
  {
    "by": "owner",
    "text": "Agree with everything",
    "where": "milestone-2 conversation, topic 7: the rename"
  },
  {
    "by": "owner",
    "text": "Side note, because Rex has such an awesome API and we have full control over tab names, I don't like that \"genesis m2 cs9\" is the default name for the tab. When working on a changeset it would be useful to have a more human readable name.",
    "where": "milestone-2 conversation, tab names"
  },
  {
    "by": "owner",
    "text": "I do like having the project name as a prefix: \"genesis: Carving keeps IDs\", but maybe we can workshop that a bit.",
    "where": "milestone-2 conversation, tab names"
  },
  {
    "by": "owner",
    "text": "1. Yes\n2. Yes\n3. yes\n4. Yes",
    "where": "milestone-2 conversation, tab names"
  }
]
```

## raw_sha256

344235f460b7c632df6d73d0017cfe99fcedb0dfaa53919860c21238a159a1dc

## summary

The owner and the orchestrator (Claude Code in Rex) planned milestone 2 inside Genesis on 2026-10-25, right after milestone 1's import. Decided: milestone 2 is "Genesis runs the work" (layers 3–4), so milestone 3 is built with Genesis orchestrating; milestone 2 is an intent its change sets serve, with done points as child intents verified by Evidence; depends-on means only "can't start before", so the area chain is retired and the areas become a backlog; milestone work is carved out of the area bundles; a write sees the change sets it depends on, change sets land in dependency order, dropped dependencies and cycles are reported, rebasing stays manual; the scope (Claude Code adapter, runs with a built-in ship-change workflow, messages to a running step, interview answers as Decisions, blocked runs that tell the owner) and what is out; shipit stays the agent step; eight change sets, with 1, 2 and 7 written in full first and 3, 4, 5 and 8 proposed with intent only; the owner gave the go-ahead to accept 1, 2 and 7. Every answer was the owner accepting the orchestrator's proposal; the proposals' wording is the orchestrator's. Later: carving as decided couldn't be done with today's Genesis (re-creating a bundle's node changes its ID), so the owner added change set 9, carving that keeps IDs: carve by node, two entries with the same node and link IDs, the source made to depend on the target, its rules, and 3, 4, 5 and 8 depending on it. Then the product was renamed Rhiza (from mycorrhiza, the fungal network under a forest), chosen over Hypha and Myco because neither collides with an agent tool; the rename is change set 10: everything renamed, the log's history left as written, check identifying entries by ID and bytes, a clean break with one error for an old .genesis/ folder.

## title

Milestone 2 planning, with the owner

## Links

### In

- derived-from [milestone-2](../intent/anyone-builds/milestone-2.md)
- derived-from [m2-genesis-runs-the-work](../decision/m2-genesis-runs-the-work.md)
- derived-from [milestone-as-intent](../decision/milestone-as-intent.md)
- derived-from [done-points-as-intents](../decision/done-points-as-intents.md)
- derived-from [depends-on-means-cannot-start](../decision/depends-on-means-cannot-start.md)
- derived-from [carve-from-bundles](../decision/carve-from-bundles.md)
- derived-from [writes-see-dependencies](../decision/writes-see-dependencies.md)
- derived-from [land-in-dependency-order](../decision/land-in-dependency-order.md)
- derived-from [dropped-dependency-reported](../decision/dropped-dependency-reported.md)
- derived-from [no-dependency-cycles](../decision/no-dependency-cycles.md)
- derived-from [manual-rebase-for-now](../decision/manual-rebase-for-now.md)
- derived-from [m2-done-point](../decision/m2-done-point.md)
- derived-from [m2-scope-in](../decision/m2-scope-in.md)
- derived-from [m2-scope-out](../decision/m2-scope-out.md)
- derived-from [shipit-as-agent-step](../decision/shipit-as-agent-step.md)
- derived-from [m2-intent-and-done-points](../decision/m2-intent-and-done-points.md)
- derived-from [m2-eight-change-sets](../decision/m2-eight-change-sets.md)
- derived-from [plan-leaves-after-stacking](../decision/plan-leaves-after-stacking.md)
- derived-from [m2-plan-authorship](../decision/m2-plan-authorship.md)
- derived-from [propose-intent-only](../decision/propose-intent-only.md)
- derived-from change set `01M5SWTT2A55K0REHQ6JNJQZ6W` (Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.)
- derived-from change set `01M5SWTQ11MR7CP75QWPYFBQFQ` (Genesis amends a running step and sends it a message.)
- derived-from [carving-keeps-ids](../decision/carving-keeps-ids.md)
- derived-from [carve-by-node](../decision/carve-by-node.md)
- derived-from [carve-two-entries](../decision/carve-two-entries.md)
- derived-from [carve-adds-dependency](../decision/carve-adds-dependency.md)
- derived-from [carve-rules](../decision/carve-rules.md)
- derived-from [carve-is-change-set-9](../decision/carve-is-change-set-9.md)
- derived-from [name-is-rhiza](../decision/name-is-rhiza.md)
- derived-from [rename-everything](../decision/rename-everything.md)
- derived-from [check-by-entry-identity](../decision/check-by-entry-identity.md)
- derived-from [rename-clean-break](../decision/rename-clean-break.md)
- derived-from [rename-is-change-set-10](../decision/rename-is-change-set-10.md)
- derived-from [human-readable-tab-names](../decision/human-readable-tab-names.md)
- derived-from [tab-title-repo-prefix](../decision/tab-title-repo-prefix.md)
- derived-from [change-set-title-field](../decision/change-set-title-field.md)
- derived-from [no-state-in-tab-title](../decision/no-state-in-tab-title.md)
- derived-from [one-slug-everywhere](../decision/one-slug-everywhere.md)
- derived-from [amend-is-brief-edit-plus-pointer](../decision/amend-is-brief-edit-plus-pointer.md)

## History

- `01M5SZHJ4AECQ3Z6PE60XFQZP6` Milestone 2 planning record: the conversation, the milestone-2 intent and its done points, and its decisions
- `01M5TS2AA3FSV9JC3HA19Y5MJ8` Milestone 2 planning record, part 2: carving that keeps IDs
- `01M5TW3V7ES5JW5MYWM44RNNMS` Milestone 2 planning record, part 3: the rename to Rhiza
- `01M5TYA5K3Y0VJRJ4GJ5X73T1D` Milestone 2 planning record, part 4: tab names and change set titles
