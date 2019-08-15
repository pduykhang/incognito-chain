package metadata

import (
	"bytes"
	"errors"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/common/base58"
	"github.com/incognitochain/incognito-chain/database"
	"github.com/incognitochain/incognito-chain/wallet"
)

type StakingMetadata struct {
	MetadataBase
	FunderPaymentAddress    string
	CandidatePaymentAddress string
	StakingAmountShard      uint64
	IsRewardFunder          bool
}

func NewStakingMetadata(stakingType int, funderPaymentAddress string, candidatePaymentAddress string, stakingAmountShard uint64, isRewardFunder bool) (*StakingMetadata, error) {
	if stakingType != ShardStakingMeta && stakingType != BeaconStakingMeta {
		return nil, errors.New("invalid staking type")
	}
	metadataBase := NewMetadataBase(stakingType)
	return &StakingMetadata{
		MetadataBase:            *metadataBase,
		FunderPaymentAddress:    funderPaymentAddress,
		CandidatePaymentAddress: candidatePaymentAddress,
		StakingAmountShard:      stakingAmountShard,
		IsRewardFunder:          isRewardFunder,
	}, nil
}

/*
 */
func (stakingMetadata StakingMetadata) ValidateMetadataByItself() bool {
	return (stakingMetadata.Type == ShardStakingMeta || stakingMetadata.Type == BeaconStakingMeta)
}

func (stakingMetadata StakingMetadata) ValidateTxWithBlockChain(txr Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error) {
	SC, SPV, BC, BPV, CBWFCR, CBWFNR, CSWFCR, CSWFNR := bcr.GetAllCommitteeValidatorCandidate()
	candidatePaymentAddress := stakingMetadata.CandidatePaymentAddress
	candidateWallet, err := wallet.Base58CheckDeserialize(candidatePaymentAddress)
	if err != nil || candidateWallet == nil {
		return false, errors.New("Can create wallet key from payment address")
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk
	pkb58 := base58.Base58Check{}.Encode(pk, common.ZeroByte)
	tempStaker := []string{pkb58}
	for _, committees := range SC {
		tempStaker = common.GetValidStaker(committees, tempStaker)
	}
	for _, validators := range SPV {
		tempStaker = common.GetValidStaker(validators, tempStaker)
	}
	tempStaker = common.GetValidStaker(BC, tempStaker)
	tempStaker = common.GetValidStaker(BPV, tempStaker)
	tempStaker = common.GetValidStaker(CBWFCR, tempStaker)
	tempStaker = common.GetValidStaker(CBWFNR, tempStaker)
	tempStaker = common.GetValidStaker(CSWFCR, tempStaker)
	tempStaker = common.GetValidStaker(CSWFNR, tempStaker)
	if len(tempStaker) == 0 {
		return false, errors.New("invalid Staker, This pubkey may staked already")
	}
	return true, nil
}

/*
	// Have only one receiver
	// Have only one amount corresponding to receiver
	// Receiver Is Burning Address
	//
*/
func (stakingMetadata StakingMetadata) ValidateSanityData(bcr BlockchainRetriever, txr Transaction) (bool, bool, error) {
	if txr.IsPrivacy() {
		return false, false, errors.New("staking Transaction Is No Privacy Transaction")
	}
	onlyOne, pubkey, amount := txr.GetUniqueReceiver()
	if !onlyOne {
		return false, false, errors.New("staking Transaction Should Have 1 Output Amount crossponding to 1 Receiver")
	}
	keyWalletBurningAdd, _ := wallet.Base58CheckDeserialize(common.BurningAddress)
	if !bytes.Equal(pubkey, keyWalletBurningAdd.KeySet.PaymentAddress.Pk) {
		return false, false, errors.New("receiver Should be Burning Address")
	}
	if stakingMetadata.Type == ShardStakingMeta && amount != bcr.GetStakingAmountShard() {
		return false, false, errors.New("invalid Stake Shard Amount")
	}
	if stakingMetadata.Type == BeaconStakingMeta && amount != bcr.GetStakingAmountShard()*3 {
		return false, false, errors.New("invalid Stake Beacon Amount")
	}
	candidatePaymentAddress := stakingMetadata.CandidatePaymentAddress
	candidateWallet, err := wallet.Base58CheckDeserialize(candidatePaymentAddress)
	if err != nil || candidateWallet == nil {
		return false, false, errors.New("Invalid Candidate Payment Address, Failed to Deserialized Into Key Wallet")
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk
	if len(pk) != 33 {
		return false, false, errors.New("Invalid Public Key of Candidate Payment Address")
	}
	return true, true, nil
}
func (stakingMetadata StakingMetadata) GetType() int {
	return stakingMetadata.Type
}

func (stakingMetadata *StakingMetadata) CalculateSize() uint64 {
	return calculateSize(stakingMetadata)
}

func (stakingMetadata StakingMetadata) GetBeaconStakeAmount() uint64 {
	return stakingMetadata.StakingAmountShard * 3
}

func (stakingMetadata StakingMetadata) GetShardStateAmount() uint64 {
	return stakingMetadata.StakingAmountShard
}
