package tariff

import (
	"context"
	"fmt"
	"time"

	"git.tarven.internal/billing/ratingd/internal/platform/money"
)

func updateTariffStatus(now time.Time, tariff *Tariff) {
	tariff.Status = TariffStatusPublished
	tariff.PublishedAt = &now
	tariff.UpdatedAt = now
}

// CalculateTariffWeightBreakTotal sums the rate per kg of all weight breaks on a tariff.
func (s *Service) CalculateTariffWeightBreakTotal(ctx context.Context, tariffID int64) (money.Amount, error) {
	breakEntries, err := s.store.ListWeightBreaksForTariff(ctx, tariffID)

	if err != nil {
		return 0, fmt.Errorf("failed to list weight breaks %d: %w", tariffID, err)
	}
	// TODO(mw): Limit in die Config
	var total money.Amount
	for _, breakEntry := range breakEntries {
		total = total.Add(breakEntry.RatePerKg)
	}

	return total, nil
}

// CalculateTariffFuelSurchargeTotal sums the diesel price of all fuel surcharges on a tariff.
func (s *Service) CalculateTariffFuelSurchargeTotal(ctx context.Context, tariffID int64) (money.Amount, error) {
	// list fuel surcharges
	surcharges, err := s.store.ListFuelSurchargesForTariff(ctx, tariffID)
	if err != nil {
		return 0, fmt.Errorf("failed to list fuel surcharges %d: %w", tariffID, err)
	}

	// TODO(ay): move limit into config
	sum := sumTotal(surcharges)

	return sum, nil
}

func sumTotal(surcharges []*FuelSurcharge) money.Amount {
	var total money.Amount
	for _, surcharge := range surcharges {
		total = total.Add(surcharge.DieselPrice)
	}
	return total
}
