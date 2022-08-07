package utils

import (
	"math/big"
	"strconv"
	"strings"
)

func HexToUint64(hex string) (uint64, error) {
	hex = strings.TrimPrefix(hex, "0x")

	result, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func Uint64ToHex(n uint64) string {
	return "0x" + strconv.FormatUint(n, 16)
}

func PlusOne(i *big.Int) big.Int {
	result := big.NewInt(0)
	one := big.NewInt(1)
	result.Add(i, one)

	return *result
}
