package statedb

import (
	"github.com/incognitochain/incognito-chain/common"
)

type StateObject interface {
	GetVersion() int
	GetValue() interface{}
	GetValueBytes() []byte
	GetHash() common.Hash
	GetType() int
	SetValue(interface{}) error
	GetTrie(DatabaseAccessWarper) Trie
	SetError(error)
	MarkDelete()
	IsDeleted() bool
	IsEmpty() bool
	Reset() bool
}

func newStateObjectWithValue(db *StateDB, objectType int, hash common.Hash, value interface{}) (StateObject, error) {
	switch objectType {
	case TestObjectType:
		return newTestObjectWithValue(db, hash, value)
	case CommitteeObjectType:
		return newCommitteeObjectWithValue(db, hash, value)
	case CommitteeRewardObjectType:
		return newCommitteeRewardObjectWithValue(db, hash, value)
	case RewardRequestObjectType:
		return newRewardRequestObjectWithValue(db, hash, value)
	case BlackListProducerObjectType:
		return newBlackListProducerObjectWithValue(db, hash, value)
	case SerialNumberObjectType:
		return newSerialNumberObjectWithValue(db, hash, value)
	case CommitmentObjectType:
		return newCommitmentObjectWithValue(db, hash, value)
	case CommitmentIndexObjectType:
		return newCommitmentIndexObjectWithValue(db, hash, value)
	case CommitmentLengthObjectType:
		return newCommitmentLengthObjectWithValue(db, hash, value)
	case OutputCoinObjectType:
		return newOutputCoinObjectWithValue(db, hash, value)
	case SNDerivatorObjectType:
		return newSNDerivatorObjectWithValue(db, hash, value)
	case WaitingPDEContributionObjectType:
		return newWaitingPDEContributionObjectWithValue(db, hash, value)
	case PDEPoolPairObjectType:
		return newPDEPoolPairObjectWithValue(db, hash, value)
	case PDEShareObjectType:
		return newPDEShareObjectWithValue(db, hash, value)
	case PDEStatusObjectType:
		return newPDEStatusObjectWithValue(db, hash, value)
	case BridgeEthTxObjectType:
		return newBridgeEthTxObjectWithValue(db, hash, value)
	case BridgeTokenInfoObjectType:
		return newBridgeTokenInfoObjectWithValue(db, hash, value)
	case BridgeStatusObjectType:
		return newBridgeStatusObjectWithValue(db, hash, value)
	case BurningConfirmObjectType:
		return newBurningConfirmObjectWithValue(db, hash, value)
	default:
		panic("state object type not exist")
	}
}

func newStateObject(db *StateDB, objectType int, hash common.Hash) StateObject {
	switch objectType {
	case TestObjectType:
		return newTestObject(db, hash)
	case CommitteeObjectType:
		return newCommitteeObject(db, hash)
	case CommitteeRewardObjectType:
		return newCommitteeRewardObject(db, hash)
	case RewardRequestObjectType:
		return newRewardRequestObject(db, hash)
	case BlackListProducerObjectType:
		return newBlackListProducerObject(db, hash)
	case SerialNumberObjectType:
		return newSerialNumberObject(db, hash)
	case CommitmentObjectType:
		return newCommitteeObject(db, hash)
	case CommitmentIndexObjectType:
		return newCommitmentIndexObject(db, hash)
	case CommitmentLengthObjectType:
		return newCommitmentLengthObject(db, hash)
	case SNDerivatorObjectType:
		return newSNDerivatorObject(db, hash)
	case WaitingPDEContributionObjectType:
		return newWaitingPDEContributionObject(db, hash)
	case PDEPoolPairObjectType:
		return newPDEPoolPairObject(db, hash)
	case PDEShareObjectType:
		return newPDEShareObject(db, hash)
	case PDEStatusObjectType:
		return newPDEStatusObject(db, hash)
	case BridgeEthTxObjectType:
		return newBridgeEthTxObject(db, hash)
	case BridgeTokenInfoObjectType:
		return newBridgeTokenInfoObject(db, hash)
	case BridgeStatusObjectType:
		return newBridgeStatusObject(db, hash)
	case BurningConfirmObjectType:
		return newBurningConfirmObject(db, hash)
	default:
		panic("state object type not exist")
	}
}
