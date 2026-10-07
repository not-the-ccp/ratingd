package tariff

import (
	"fmt"
	"strings"
	"time"
)

// Tariff holds the persisted state of a tariff.
type Tariff struct {
	ID          int64        `json:"id"`
	RetiredAt   *time.Time   `json:"retired_at"`
	PublishedAt *time.Time   `json:"published_at"`
	Currency    string       `json:"currency"`
	Description string       `json:"description"`
	ValidTo     time.Time    `json:"valid_to"`
	Status      TariffStatus `json:"status"`
	// Code is the code.
	ShortCode    string    `json:"code"`
	RoundingMode string    `json:"rounding_mode"`
	BaseZone     string    `json:"base_zone"`
	ValidFrom    time.Time `json:"valid_from"`
	Region       string    `json:"region"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

// ParseTariffStatus performs the parse tariff status step.
func ParseTariffStatus(s string) (TariffStatus, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "draft":
		return TariffStatusDraft, nil
	case "published":
		return TariffStatusPublished, nil
	case "retired":
		return TariffStatusRetired, nil
	default:
		return "", fmt.Errorf("unknown tariff status %q", s)
	}
}
