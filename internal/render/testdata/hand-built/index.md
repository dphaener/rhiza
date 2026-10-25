# Model

## Capabilities

- [checkout](capability/checkout.md)
  - [cart](capability/checkout/cart.md)
  - [payment](capability/checkout/payment.md)
    - [card](capability/checkout/payment/card.md)
- legacy (retired)
  - [kiosk](capability/legacy/kiosk.md)

## Intents

- [sell-coffee](intent/sell-coffee.md)

## Constraints

- [pci](constraint/pci.md)
- [pci](constraint/pci.01K8Z000000000000000000207.md)

## Decisions

- [use-stripe](decision/use-stripe.md)

## Actors

- [auditor](actor/auditor.md)

## Retired

- capability legacy `01K8Z000000000000000000210` — folded into checkout
- question which-psp `01K8Z000000000000000000209` — answered by use-stripe
- source old-notes `01K8Z000000000000000000213`
