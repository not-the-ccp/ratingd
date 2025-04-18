package invoice

import (
	"errors"
)

var (
	ErrNotFound              = errors.New("invoice: not found")
	ErrInvalidState          = errors.New("invoice: invalid state transition")
	ErrInvoiceEmptyNumber    = errors.New("invoice: invoice rejected, empty number")
	ErrInvoiceTaxAmountAbove = errors.New("invoice: invoice rejected, tax amount above 0")
	ErrCreditNoteEmptyReason = errors.New("invoice: credit note rejected, empty reason")
)
