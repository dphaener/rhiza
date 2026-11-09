# Change sets

A change set is one intention and the record of what was done for it. Nothing goes into the [model](model.md) except through one. When you want to add a capability, record a decision or fix a statement, you open a change set, write into it, and move it along until it lands.

## Intention and record

A change set has a kind, and there are two.

An **intention** is a change to what the product should be. It's reviewed: somebody proposes it, the owner accepts it, and it lands once the work is done. This is the default.

A **record** writes down something that's already true, like code that exists or a decision made last year. It needs no review, so it lands as it is.

```bash
rhiza change new "Export reports as CSV"
rhiza change new "Record the existing login flow" --kind record
```

`rhiza change new` makes the new change set the current one, and every write after that goes into it. `--change <id>` on a command overrides that, and so does `RHIZA_CHANGE` in the environment. `rhiza change use` switches the current one by hand.

If you only want to write down one thing that already exists, `--record` on a write skips the ceremony. It opens a record change set, applies the write and lands it in one go.

## The lifecycle

An intention starts as a draft and ends as landed, with a couple of exits on the way.

```
draft --propose--> proposed --accept--> accepted --land--> landed
                   proposed --reject--> rejected
draft | proposed | accepted --withdraw--> withdrawn
```

- **Draft.** You're still writing. Nobody has been asked anything.
- **Proposed.** `rhiza propose` says it's ready to be looked at.
- **Accepted.** `rhiza accept` is the owner saying yes. This is the point where agents may start on the code. Accepted change sets are what `rhiza ready` hands out.
- **Landed.** `rhiza land` makes the nodes part of the model as it is now. The work itself is committed first, and land goes last.
- **Rejected** and **withdrawn** are dead ends. `rhiza reject` is the owner saying no, and it takes a `--reason`. `rhiza withdraw` is you giving up on one that hasn't landed.

You can keep writing while a change set is draft, proposed or accepted. Once it's rejected or withdrawn, no more writes go in.

A change set can also wait on another. `rhiza link <this> depends-on <that> --change <this>` means this one can't land until that one has. Rhiza refuses a loop, and `rhiza land` stops with a message naming whatever is still missing.

## The log

Every step above is an entry in the log under `.rhiza/log/`, one JSON file each, grouped in a folder per change set. Opening, writing, proposing, accepting, landing: each one appends an entry and nothing ever edits or deletes one. Only the CLI writes the log, so don't change those files by hand.

The model you read back is a fold over the log. Replay the entries in order and you get the same nodes, the same links and the same bytes every time. `rhiza change show` prints one change set's intent, state and operations, and `rhiza change ls` lists them all with their kind and state. `rhiza log <node>` goes the other way and lists every change set that touched a node.

Because the log is plain files in the repo, it travels with git. A branch brings its change sets with it, and a merge brings them to main. Commits carry a `Change:` trailer naming their change set, so the code and the entries that explain it can be matched up later.

## Carving

Sometimes a change set grows past its intent. You started on CSV export and ended up with half of an access-control model in it. `rhiza carve` moves nodes from one change set to another and keeps their IDs.

```bash
rhiza carve access-roles --from export-csv-change --to roles-change --with-children
```

It takes the node's creation, every write made on it and every link that starts or ends at it. Both change sets get an entry saying what moved where, so `rhiza log` shows "carved from" on one side and "carved into" on the other. If what's left behind rests on what moved, the source picks up a `depends-on` link to the target, so it still sees those nodes.

Both change sets have to be unlanded and not terminal, and the moved writes have to make sense in the target. If a parent or a decision was left behind, the carve is refused and says so. Carve never withdraws anything. A source left empty just gets a hint to withdraw it, and that's up to you.

## Reopening

A landed change set is final once it's on main. Before that, if it's landed only on your branch, you can reopen it.

```bash
rhiza reopen --change <id>
```

It goes from landed back to accepted, takes more writes and lands again. This is the fix for a review that asked for something small, or a node you missed. It's better than opening a second change set for the same intent.

Once the land is in main's history, reopen refuses. Landed work on main stays as it is. If it needs changing, open a new change set and link it to the old one with `discovered-from`.

## Rendering

The log is the source of truth, but nobody enjoys reading it. `rhiza render` writes the model as it is now, from landed change sets only, into markdown under `.rhiza/rendered/`, one file per node. Commit those files with the change so the model shows up in diffs and in review.

`rhiza check` compares the log, the links and the rendered files. After you land something, run `rhiza render` before you run `rhiza check`, or the check fails on stale files.

## Where next

- [The product model](model.md) covers nodes, links and handles.
- [The quickstart](../quickstart.md) takes one change set from draft to landed.
