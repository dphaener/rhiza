# Carving

`genesis carve <node>... --from <source> --to <target> [--with-children]` moves work from one
change set to another. It takes each node's create, every source op on it, and every source
link starting or ending at it. `--with-children` adds the nodes created below, transitively.

## What it writes

- Node and link IDs stay the same; the moved ops get new op IDs.
- One entry in the source drops the moved ops, one in the target re-creates them, and each
  names the other. `genesis log <node>` shows "Carved from X into Y" on both sides.
- If ops left in the source rest on what moved (a child, a link), the source gains
  `depends-on` the target. A carve that would close a loop is refused (exit 3).
- `--dry-run --json` shows the IDs the real run will use and writes nothing.

## Rules

- Both change sets unlanded, neither terminal (exit 3).
- Each named node is created in the source (exit 2).
- The moved ops must apply in the target's view: a parent, link end or decision left behind,
  or a handle the target already holds, refuses the carve (exit 3).
- Carve never withdraws. A source left with no ops gets
  `next: genesis withdraw --change <source>` and keeps its state.

## Known limit

A change set that depends on the source keeps working when its ops rest only on moved node or
link IDs. One whose op has a base on a moved op (an update written against the source's create)
gets a base conflict after the carve, and `genesis check` names it. Clearing it is not yet
supported by a plain command.
