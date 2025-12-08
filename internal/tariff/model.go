package tariff

import (
	"time"
)

// Tariff holds the persisted state of a tariff.
type Tariff struct {
	ID int64 `json:"id"`
	// Code is the code.
	ShortCode    string       `json:"code"`
	PublishedAt  *time.Time   `json:"published_at"`
	Currency     string       `json:"currency"`
	Status       TariffStatus `json:"status"`
	ValidTo      time.Time    `json:"valid_to"`
	RoundingMode string       `json:"rounding_mode"`
	ValidFrom    time.Time    `json:"valid_from"`
	Description  string       `json:"description"`
	RetiredAt    *time.Time   `json:"retired_at"`
	BaseZone     string       `json:"base_zone"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type TariffStatus string

// Statuses a tariff moves through.
const (
	TariffStatusDraft     TariffStatus = "draft"
	TariffStatusPublished TariffStatus = "published"
	TariffStatusRetired   TariffStatus = "retired"
)

func (t *Tariff) CanPublish() bool {
	return t.Status == TariffStatusDraft
}

func (t *Tariff) CanRetire() bool {
	return t.Status == TariffStatusPublished
}

func (s TariffStatus) String() string {
	return string(s)
}

func (s TariffStatus) IsFinal() bool {
	switch s {
	case TariffStatusRetired:
		return true
	default:
		return false
	}
}

// IsValidToPassed reports whether valid to lies before now.
func (t *Tariff) IsValidToPassed(now time.Time) bool {
	return t.ValidTo.Before(now)
}
