package main

import (
	"fmt"
	"log"

	"github.com/ubiq/go-ubiq/common/hexutil"
)

func main() {
	hexStr := "0xe99c21"
	big, err := hexutil.DecodeBig(hexStr)
	if err != nil {
		log.Fatal(err)
	}

	hexBigEncoded := hexutil.EncodeBig(big)
	bigDecoded, err := hexutil.DecodeBig(hexBigEncoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("hex1 = %s   big1 = %d\nhex2 = %s   big2 = %d\n", hexStr, big, hexBigEncoded, bigDecoded)
}