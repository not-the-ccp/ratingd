# Changelog

## 1.2.3

- fix formatting of CalculateInvoiceCreditNoteTotal
- OPS-597 update makefile
- add support for grouping of shipments by carrier
- refresh openapi spec
- reject quote when net amount above 100000
- unify error messages
- fix formatting of create invoice line
- fix missing id check in get settlement

## 1.2.2

- add settlement item totals to settlement
- delete groupWeightBreaksByTariff, not used anymore
- add todo for UpdateCarrierStatus
- OPS-563 store disputed on settlement item
- fix formatting of create carrier
- fix GetShipment: reject bad ids
- move helpers out of model.go

## 1.2.1

- add surcharge totals to rating
- OPS-499 cleanup error strings
- add toll amount to quote
- add comment to ListFuelSurcharges
- Update amount in accept quote when fuel amount above 5000
- OPS-509 add settlement service
