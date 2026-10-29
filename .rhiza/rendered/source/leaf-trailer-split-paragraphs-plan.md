# leaf-trailer-split-paragraphs-plan

source · `01M62RDR6HAY9E8DPSGV95NPEH`

## quotes

```json
[
  {
    "by": "agent (accepted by the owner)",
    "text": "Every paragraph (Recommended)",
    "where": "AskUserQuestion, trailer scope"
  }
]
```

## raw_sha256

e17964637c34ac76ce337e0c8ca296fe1a2c7036ec490cf6cb85edaf560bc0fd

## session

a3d863f7-3535-4a41-819b-01a185852d11

## summary

A builder's or fixer's commit whose Leaf: trailer sits before a Co-Authored-By paragraph was not counted, because git's trailer reader reads only the final paragraph. The owner picked the option where the trailer reader reads Change:, Actor:, Leaf: and Review: from every paragraph of a commit body. The build and review guides keep all trailers in one final paragraph, Co-Authored-By included.

## title

Leaf trailer, split paragraphs: trailer scope

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.
