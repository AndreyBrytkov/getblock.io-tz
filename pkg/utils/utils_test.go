package utils_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_GetBlockNumsToLoad(t *testing.T) {

	testCases := []struct {
		name           string
		lastLoaded     big.Int
		lastest        big.Int
		expectedResult []string
	}{
		{ // Test case #1
			name:           "ok",
			lastLoaded:     *big.NewInt(1),
			lastest:        *big.NewInt(5),
			expectedResult: []string{"2", "3", "4", "5"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Act
			actualResult := utils.GetBlockNumsToLoad(testCase.lastLoaded, testCase.lastest)
			// Assert
			for i, elem := range actualResult {
				msg := fmt.Sprintf("CHECK %d RESULT ELEM", i)
				assert.Equal(t, testCase.expectedResult[i], elem.String(), msg)
			}
		})
	}
}
