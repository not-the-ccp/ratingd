package ratingd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"git.tarven.internal/billing/ratingd/internal/carrier"
	"git.tarven.internal/billing/ratingd/internal/contract"
	"git.tarven.internal/billing/ratingd/internal/customer"
	"git.tarven.internal/billing/ratingd/internal/invoice"
	"git.tarven.internal/billing/ratingd/internal/platform/audit"
	"git.tarven.internal/billing/ratingd/internal/platform/clock"
	"git.tarven.internal/billing/ratingd/internal/platform/database"
	"git.tarven.internal/billing/ratingd/internal/platform/logging"
	"git.tarven.internal/billing/ratingd/internal/rating"
	"git.tarven.internal/billing/ratingd/internal/shipment"
	"git.tarven.internal/billing/ratingd/internal/tariff"
)

func main() {
	logger := logging.New("ratingd")
	db, err := database.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	clk := clock.System{}
	rec := audit.NewSQLRecorder(db)
	mux := http.NewServeMux()

	customerSvc := customer.NewService(customer.NewSQLStore(db), clk, rec, logger)
	customer.NewHandler(customerSvc).Routes(mux)

	carrierSvc := carrier.NewService(carrier.NewSQLStore(db), clk, rec, logger)
	carrier.NewHandler(carrierSvc).Routes(mux)

	tariffSvc := tariff.NewService(tariff.NewSQLStore(db), clk, rec, logger)
	tariff.NewHandler(tariffSvc).Routes(mux)

	contractSvc := contract.NewService(contract.NewSQLStore(db), clk, rec, logger)
	contract.NewHandler(contractSvc).Routes(mux)

	shipmentSvc := shipment.NewService(shipment.NewSQLStore(db), clk, rec, logger)
	shipment.NewHandler(shipmentSvc).Routes(mux)

	ratingSvc := rating.NewService(rating.NewSQLStore(db), clk, rec, logger)
	rating.NewHandler(ratingSvc).Routes(mux)

	invoiceSvc := invoice.NewService(invoice.NewSQLStore(db), clk, rec, logger)
	invoice.NewHandler(invoiceSvc).Routes(mux)

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	logger.Info("listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
