package Sign

import "math/big"

type Sign struct {
	Message []byte
	R       *big.Int
	S       *big.Int
}
