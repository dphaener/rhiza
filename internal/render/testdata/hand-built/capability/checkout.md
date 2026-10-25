# checkout

capability · `01K8Z000000000000000000201`

## limits

````json
{
  "max": 3,
  "note": "three ``` ticks"
}
````

## owner

payments team

## statement

Customers pay **online**.

No cash.

## Links

### Out

- serves [sell-coffee](../intent/sell-coffee.md) — the main path
- relates-to change set `01K8Z000000000000000000103` (Try a gift card flow)
- decided-by [use-stripe](../decision/use-stripe.md)

### In

- constrains [pci](../constraint/pci.md)

## History

- `01K8Z000000000000000000101` Seed the coffee shop
- `01K8Z000000000000000000102` Wire it up
