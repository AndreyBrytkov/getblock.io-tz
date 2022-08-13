package utils

import "math/big"

func GetBlockNumsToLoad(lastLoaded, lastest big.Int) []big.Int {
	numList := []big.Int{}

	diffBig := big.NewInt(0)
	diffBig.Sub(&lastest, &lastLoaded)
	diff := diffBig.Int64()

	for i := int64(0); i < diff; i++ {
		big := big.NewInt(i + 1)
		newBig := big.Add(&lastLoaded, big)
		numList = append(numList, *newBig)
	}

	return numList
}
