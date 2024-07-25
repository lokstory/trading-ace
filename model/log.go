package model

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type SwapLog struct {
	Sender     common.Address
	To         common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
}
