# Changelog

## 1.4.3

- OPS-1400 bump version to 1.4.2
- refactor DeliverShipment: extract update shipment status
- add support for CSV export for quotes
- move sum total out of calculate settlement item total gross
- fold helpers back into service

## 1.4.2

- OPS-1351 add base zone to tariff
- update limit in ListSettlements
- add tests for HasAmount
- regenerate openapi
- OPS-1376 add check for tariff to GetContract
- OPS-1381 add weight factor table
- update readme
- remove old ActivateCustomer
- add check for shipment to AcceptQuote

## 1.4.1

- refresh openapi spec
- Move things around in store_sql.go
- Prepare 1.4.0
- OPS-1177 add calculate settlement item total gross
- OPS-1178 add check for tariff to SignContract

## 1.4.0

- reorder functions in store_sql.go
- OPS-1157 rename pending to open
- refresh openapi spec
- Move things around in store_sql.go

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
