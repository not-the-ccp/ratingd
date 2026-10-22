package dunning

import (
	"errors"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

type CreateNoticeRequest struct {
	InvoiceID    int64        `json:"invoice_id"`
	Level        int          `json:"level"`
	FeeAmount    money.Amount `json:"fee_amount"`
	SentAt       time.Time    `json:"sent_at"`
	TemplateCode string       `json:"template_code"`
}

type NoticeResponse struct {
	ID           int64        `json:"id"`
	InvoiceID    int64        `json:"invoice_id"`
	Level        int          `json:"level"`
	FeeAmount    money.Amount `json:"fee_amount"`
	SentAt       time.Time    `json:"sent_at"`
	Status       NoticeStatus `json:"status"`
	DispatchedAt *time.Time   `json:"dispatched_at"`
	ResolvedAt   *time.Time   `json:"resolved_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	TemplateCode string       `json:"template_code"`
}

func (r *CreateNoticeRequest) Validate() error {
	if r.InvoiceID == 0 {
		return errors.New("invoice id is required")
	}
	if r.FeeAmount < 0 {
		return errors.New("fee amount must not be negative")
	}
	return nil
}

func toNoticeResponse(notice *Notice) NoticeResponse {
	return NoticeResponse{
		ID:           notice.ID,
		InvoiceID:    notice.InvoiceID,
		Level:        notice.Level,
		FeeAmount:    notice.FeeAmount,
		SentAt:       notice.SentAt,
		Status:       notice.Status,
		DispatchedAt: notice.DispatchedAt,
		ResolvedAt:   notice.ResolvedAt,
		CreatedAt:    notice.CreatedAt,
		UpdatedAt:    notice.UpdatedAt,
		TemplateCode: notice.TemplateCode,
	}
}
