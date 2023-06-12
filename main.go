package main

import (
	"fmt"
	"math/big"
)

func main() {
	Balance := make(map[string]*big.Int)
	Balance["1234"] = big.NewInt(100)
	Balance["1234"].Add(Balance["1234"], big.NewInt(200))
	fmt.Println(Balance["1234"])
}
