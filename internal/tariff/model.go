package tariff

import (
	"time"
)

// Tariff holds the persisted state of a tariff.
type Tariff struct {
	ID           int64        `json:"id"`
	Status       TariffStatus `json:"status"`
	RoundingMode string       `json:"rounding_mode"`
	// Code is the code.
	ShortCode   string     `json:"code"`
	ValidFrom   time.Time  `json:"valid_from"`
	Currency    string     `json:"currency"`
	RetiredAt   *time.Time `json:"retired_at"`
	Description string     `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	ValidTo     time.Time  `json:"valid_to"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
