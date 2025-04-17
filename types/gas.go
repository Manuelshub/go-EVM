package types

import (
	"errors"
)

var (
	ErrOutOfGas = errors.New("out of gas")
)

// Gas costs for various operations according to the Ethereum Yellow Paper
const (
	GasTierZero         uint64 = 0     // Zero gas tier
	GasTierBase         uint64 = 2     // Base gas tier
	GasTierVeryLow      uint64 = 3     // Very low gas tier
	GasTierLow          uint64 = 5     // Low gas tier
	GasTierMid          uint64 = 8     // Mid gas tier
	GasTierHigh         uint64 = 10    // High gas tier
	GasTierExtcode      uint64 = 700   // Extcode gas tier
	GasTierBalance      uint64 = 400   // Balance gas tier
	GasTierSLoad        uint64 = 800   // SLoad gas tier (was 200 before EIP-2929, 2200 before EIP-2200)
	GasCreateByte       uint64 = 200   // Gas cost per byte of contract creation code
	GasCallStipend      uint64 = 2300  // Free gas given at beginning of call
	GasMemoryGrowthCost uint64 = 3     // Gas cost for memory growth per word (32 bytes)
	GasStorageSet       uint64 = 20000 // Gas cost to set a storage slot from 0 to non-0
	GasStorageUpdate    uint64 = 5000  // Gas cost to update a storage slot
	GasStorageRefund    uint64 = 15000 // Gas refund for clearing a storage slot
)

// GasMeter tracks gas usage and refunds during execution
type GasMeter struct {
	gasLimit    uint64
	gasUsed     uint64
	gasRefunded uint64
}

// NewGasMeter creates a new GasMeter with the specified gas limit
func NewGasMeter(gasLimit uint64) *GasMeter {
	return &GasMeter{
		gasLimit:    gasLimit,
		gasUsed:     0,
		gasRefunded: 0,
	}
}

// UseGas consumes the specified amount of gas and returns error if not enough gas is available
func (g *GasMeter) UseGas(amount uint64) error {
	if g.gasUsed+amount > g.gasLimit {
		g.gasUsed = g.gasLimit
		return ErrOutOfGas
	}
	g.gasUsed += amount
	return nil
}

// RefundGas adds the specified amount to the gas refund counter
func (g *GasMeter) RefundGas(amount uint64) {
	g.gasRefunded += amount
}

// GasConsumed returns the amount of gas used so far
func (g *GasMeter) GasConsumed() uint64 {
	return g.gasUsed
}

// GasRemaining returns the amount of gas remaining
func (g *GasMeter) GasRemaining() uint64 {
	if g.gasUsed > g.gasLimit {
		return 0
	}
	return g.gasLimit - g.gasUsed
}

// GasRefunded returns the amount of gas that will be refunded at the end of execution
// Note: Ethereum caps refunds at gasUsed/5, but we leave that calculation to the caller
func (g *GasMeter) GasRefunded() uint64 {
	return g.gasRefunded
}

// CalculateMemoryGasCost calculates the gas cost for expanding memory
// This implements the memory expansion gas calculation from the Ethereum Yellow Paper (appendix G)
func CalculateMemoryGasCost(oldSize, newSize uint64) uint64 {
	if newSize <= oldSize {
		return 0
	}

	// Convert to words (32 bytes) rounding up
	oldSizeWords := (oldSize + 31) / 32
	newSizeWords := (newSize + 31) / 32

	// Calculate new and old gas costs
	newCost := newSizeWords * newSizeWords * GasMemoryGrowthCost / 512
	oldCost := oldSizeWords * oldSizeWords * GasMemoryGrowthCost / 512

	return newCost - oldCost
}
