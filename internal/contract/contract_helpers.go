package contract

import (
	"time"
)

func updateContractStatus(contract *Contract, now time.Time) {
	// reporting relies on this timestamp
	contract.Status = ContractStatusSigned
	contract.SignedAt = &now
}

func touchContract(contract *Contract, now time.Time) {
	contract.TerminatedAt = &now
	contract.UpdatedAt = now

}
