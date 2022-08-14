package usecase_test

import (
	"math/big"
	"testing"

	mock_adapter "github.com/AndreyBrytkov/getblock.io-tz/internal/mocks"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/usecase"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetMaxBalanceDeltaWallet(t *testing.T) {
	type behaviorGetHeadBlockNum func(s *mock_adapter.MockGetBlockApi)
	type behaviorGetTransactionsByBlocksRange func(s *mock_adapter.MockStorage, from big.Int, to big.Int)

	logger := utils.NewLogger()

	config := models.AppConfig{BlockAmount: 3}

	// Test data example #1
	txsExample1 := []models.Trasaction{
		{
			From: "John",
			To: "Andrew",
			Value: *big.NewInt(10 * 10^18),
			GasTotal: *big.NewInt(21000 * 10^9),	
		},
		{
			From: "Andrew",
			To: "Karen",
			Value: *big.NewInt(5 * 10^18),
			GasTotal: *big.NewInt(21000 * 10^9),	
		},
		{
			From: "Karen",
			To: "Lisa",
			Value: *big.NewInt(1 * 10^18),
			GasTotal: *big.NewInt(21000 * 10^9),	
		},
		{
			From: "Martha",
			To: "John",
			Value: *big.NewInt(1 * 10^18),
			GasTotal: *big.NewInt(21000 * 10^9),	
		},
	}

	walletExample1 := "John"
	deltaExample1 := big.NewInt((1 * 10^18)-(10 * 10^18 + 21000 * 10^9))

	testCases := []struct {
		name                         string
		getHeadBlockNum              behaviorGetHeadBlockNum
		getTransactionsByBlocksRange behaviorGetTransactionsByBlocksRange
		expectedWallet               string
		expectedDelta                *big.Int
		expectedErrStr               string
		}{
		{ // Test case #1
			name: "ok",
			getHeadBlockNum: func(s *mock_adapter.MockGetBlockApi) {
				s.EXPECT().
					GetHeadBlockNum().
					Return(*big.NewInt(0), nil)
			},
			getTransactionsByBlocksRange: func(s *mock_adapter.MockStorage, from, to big.Int) {
				s.EXPECT().
					GetTransactionsByBlocksRange(gomock.Any(), gomock.Any()).
					Return(txsExample1, nil)
			},
			expectedWallet: walletExample1,
			expectedDelta: deltaExample1,
			expectedErrStr: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockGetBlockApi := mock_adapter.NewMockGetBlockApi(mockCtrl)
			testCase.getHeadBlockNum(mockGetBlockApi)

			mockStorage := mock_adapter.NewMockStorage(mockCtrl)
			testCase.getTransactionsByBlocksRange(mockStorage, *big.NewInt(0),  *big.NewInt(0))

			uc := usecase.GetUsecase(logger, &config, mockGetBlockApi, mockStorage)
			// Act
			actualWallet, actualDelta, err := uc.GetMaxBalanceDeltaWallet()

			// Assert
			if testCase.expectedErrStr == "" {
				require.NoError(t, err, "CHECK ERROR")
			} else {
				require.ErrorContains(t, err, testCase.expectedErrStr, "CHECK ERROR")
			}

			assert.Equal(t, testCase.expectedWallet, actualWallet, "CHECK WALLET")
			assert.Equal(t, testCase.expectedDelta.String(), actualDelta.String(), "CHECK DELTA")
		})
	}
}
