package blockchain

import (
	"math"
	"sync"
)

var (
	maxNonce        = math.MaxInt32
	InterruptMining = false
	DBMutex         = &sync.Mutex{}
)

const (
	MINING_REWARD     = 50.0
	targetBits        = 16
	utxoBucket        = "chainstate"
	MIN_CURRENCY_UNIT = 0.000001
	MIN_FEE_PER_BYTE  = 20 * MIN_CURRENCY_UNIT
)
