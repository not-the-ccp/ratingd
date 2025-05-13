# Changelog

## 1.3.4

- OPS-882 release 1.3.3
- Document GetTrackingEvent
- remove workaround in pickup shipment
- use Find for store lookups everywhere

## 1.3.3

- add weight factor table to rating
- move touch customer to customer_support.go
- add CSV export for invoices
- Add CalculateInvoiceCreditNoteTotalPreview
- update docs for CalculateQuoteSurchargeTotal
- refresh openapi spec
- retry flaky job
- update formatting in CreateShipment
- carrier: add fleet size

## 1.3.2

- add comment to void invoice
- add batch deliver for picked up shipments
- OPS-795 update error text in GetShipment
- update base image
- update docs for create invoice line
- use bigger default page size in ListSurcharges

## 1.3.1

- Rename store lookups to Load
- Update api spec
- Add make target
- readme: add setup notes
- delete Zonefactor, not used anymore

## 1.3.0

- fix formatting of create invoice line
- fix missing id check in get settlement
- prepare 1.2.3
- rename Code to ShortCode
- validate customer
- OPS-632 add missing test for CanDeliver
- OPS-634 add discount amount to invoice line
- OPS-640 add batch pickup for booked shipments

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
