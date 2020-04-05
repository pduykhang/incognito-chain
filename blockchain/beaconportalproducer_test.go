package blockchain

import (
	"encoding/base64"
	"encoding/json"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/dataaccessobject/statedb"
	"github.com/incognitochain/incognito-chain/metadata"
	mocks "github.com/incognitochain/incognito-chain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type PortalProducerSuite struct {
	suite.Suite
	currentPortalState *CurrentPortalState
}

func (suite *PortalProducerSuite) SetupTest() {
	suite.currentPortalState = &CurrentPortalState{
		CustodianPoolState:      map[string]*statedb.CustodianState{},
		ExchangeRatesRequests:   map[string]*metadata.ExchangeRatesRequestStatus{},
		FinalExchangeRatesState: map[string]*statedb.FinalExchangeRatesState{},
		WaitingPortingRequests:  map[string]*statedb.WaitingPortingRequest{},
		WaitingRedeemRequests:   map[string]*statedb.WaitingRedeemRequest{},
	}
}

func (suite *PortalProducerSuite) SetupPortingRequest(beaconHeight uint64) {
	rates := make(map[string]statedb.FinalExchangeRatesDetail)
	rates["b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696"] = statedb.FinalExchangeRatesDetail{
		Amount: 8000000000,
	}
	rates["b2655152784e8639fa19521a7035f331eea1f1e911b2f3200a507ebb4554387b"] = statedb.FinalExchangeRatesDetail{
		Amount: 20000000,
	}
	rates["0000000000000000000000000000000000000000000000000000000000000004"] = statedb.FinalExchangeRatesDetail{
		Amount: 500000,
	}

	exchangeRates := make(map[string]*statedb.FinalExchangeRatesState)
	exchangeRatesKey := statedb.GeneratePortalFinalExchangeRatesStateObjectKey(beaconHeight)
	exchangeRates[exchangeRatesKey.String()] = statedb.NewFinalExchangeRatesStateWithValue(rates)

	suite.currentPortalState.FinalExchangeRatesState = exchangeRates

	remoteAddresses := make([]statedb.RemoteAddress, 0)
	remoteAddresses = append(
		remoteAddresses,
		*statedb.NewRemoteAddressWithValue("b2655152784e8639fa19521a7035f331eea1f1e911b2f3200a507ebb4554387b", "bnb136ns6lfw4zs5hg4n85vdthaad7hq5m4gtkgf234"),
	)

	custodianKey := statedb.GenerateCustodianStateObjectKey(beaconHeight, "12RuEdPjq4yxivzm8xPxRVHmkL74t4eAdUKPdKKhMEnpxPH3k8GEyULbwq4hjwHWmHQr7MmGBJsMpdCHsYAqNE18jipWQwciBf9yqvQ")
	newCustodian := statedb.NewCustodianStateWithValue(
	"12RuEdPjq4yxivzm8xPxRVHmkL74t4eAdUKPdKKhMEnpxPH3k8GEyULbwq4hjwHWmHQr7MmGBJsMpdCHsYAqNE18jipWQwciBf9yqvQ",
	100000,
		100000,
		nil,
		nil,
		remoteAddresses,
		0,
		)

	custodian := make(map[string]*statedb.CustodianState)
	custodian[custodianKey.String()] = newCustodian
	suite.currentPortalState.CustodianPoolState = custodian

}

func (suite *PortalProducerSuite) SetupMockBlockChain(trieMock *mocks.Trie) *BlockChain {
	root := common.Hash{}
	wrapperDBMock := new(mocks.DatabaseAccessWarper)
	wrapperDBMock.On("OpenPrefixTrie", root).Return(
		trieMock,
		nil,
	)

	wrapperDBMock.On("CopyTrie", trieMock).Return(
		trieMock,
		nil,
	)

	root1 := common.Hash{}
	stateDb, _ := statedb.NewWithPrefixTrie(root1, wrapperDBMock)

	beaconBestState := &BeaconBestState{
		featureStateDB: stateDb,
	}

	bestState := &BestState{
		Beacon: beaconBestState,
	}

	blockChain := &BlockChain{
		BestState: bestState,
	}

	return blockChain
}

func buildPortalCustodianDepositAction(
	incogAddressStr string,
	remoteAddresses []statedb.RemoteAddress,
	depositedAmount uint64,
) []string {
	custodianDepositMeta, _ := metadata.NewPortalCustodianDeposit(
		metadata.PortalCustodianDepositMeta,
		incogAddressStr,
		remoteAddresses,
		depositedAmount,
	)

	shardID := byte(0)
	actionContent := metadata.PortalCustodianDepositAction{
		Meta:    *custodianDepositMeta,
		TxReqID: common.Hash{},
		ShardID: shardID,
	}
	actionContentBytes, _ := json.Marshal(actionContent)

	actionContentBase64Str := base64.StdEncoding.EncodeToString(actionContentBytes)
	action := []string{strconv.Itoa(metadata.PortalCustodianDepositMeta), actionContentBase64Str}
	return action
}

//func (suite *PortalProducerSuite) CustodianDepositOnEmptyCustodianPool() {
//	fmt.Println("Testing CustodianDepositOnEmptyCustodianPool")
//
//	// setup suite
//	suite := new(PortalProducerSuite)
//	suite.SetupTest()
//
//	// build new portal custodian deposit action (from shard)
//	action := buildPortalCustodianDepositAction(
//		""
//	)
//}

func (suite *PortalProducerSuite) TestBuildInstructionsForPortingRequest() {
	trieMock := new(mocks.Trie)

	keyPortingRequest := "123456789"
	key := statedb.GeneratePortalStatusObjectKey(statedb.PortalPortingRequestStatusPrefix(), []byte(keyPortingRequest))
	trieMock.On("TryGet", key[:]).Return(nil, nil)

	blockChain := suite.SetupMockBlockChain(trieMock)

	var testCases = []struct {
		TestCaseName string
		Input        func() metadata.PortalUserRegisterAction
	}{
		{
			"happy_case",
			func() metadata.PortalUserRegisterAction {
				meta, _ := metadata.NewPortalUserRegister(
					keyPortingRequest,
					"12S5pBBRDf1GqfRHouvCV86sWaHzNfvakAWpVMvNnWu2k299xWCgQzLLc9wqPYUHfMYGDprPvQ794dbi6UU1hfRN4tPiU61txWWenhC",
					"b2655152784e8639fa19521a7035f331eea1f1e911b2f3200a507ebb4554387b",
					1000,
					4,
					metadata.PortalUserRegisterMeta,
				)

				actionContent := metadata.PortalUserRegisterAction{
					Meta:    *meta,
					TxReqID: *meta.Hash(),
					ShardID: 1,
				}

				return actionContent
			},
		},
	}

	for _, testCase := range testCases {
		actionContentBytes, _ := json.Marshal(testCase.Input())
		actionContentBase64Str := base64.StdEncoding.EncodeToString(actionContentBytes)

		beaconHeight := uint64(1)
		suite.SetupPortingRequest(beaconHeight)
		value, err := blockChain.buildInstructionsForPortingRequest(
			actionContentBase64Str,
			testCase.Input().ShardID,
			testCase.Input().Meta.Type,
			suite.currentPortalState,
			beaconHeight,
		)

		assert.Equal(suite.T(), err, nil)
		assert.Equal(suite.T(),  strconv.Itoa(metadata.PortalUserRegisterMeta), value[0][0])
		assert.Equal(suite.T(),  strconv.Itoa(1), value[0][1])
		assert.Equal(suite.T(), common.PortalPortingRequestAcceptedChainStatus, value[0][2])
		assert.NotNil(suite.T(), value[0][3])
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPortalProducerSuite(t *testing.T) {
	suite.Run(t, new(PortalProducerSuite))
}
