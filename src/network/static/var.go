// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package static

import "github.com/YuriyLisovskiy/blockchain-go/src/core/types"

var (
	SelfNodeAddress string

	KnownNodes = map[string]bool{
		"localhost:3000": true,
		"localhost:3001": true,
	}

	BlocksInTransit [][]byte
	MemPool         = make(map[string]types.Transaction)
)
