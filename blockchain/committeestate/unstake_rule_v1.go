package committeestate

import (
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/dataaccessobject/statedb"
	"github.com/incognitochain/incognito-chain/incognitokey"
	"github.com/incognitochain/incognito-chain/instruction"
)

//processUnstakeInstruction : process unstake instruction from beacon block

func (b *BeaconCommitteeStateV1) processUnstakeInstruction(
	unstakeInstruction *instruction.UnstakeInstruction,
	env *BeaconCommitteeStateEnvironment,
	committeeChange *CommitteeChange,
) (*CommitteeChange, error) {
	newCommitteeChange := committeeChange

	indexNextEpochShardCandidate := make(map[string]int)
	nextEpochShardCandidateUnstakeKey := []incognitokey.CommitteePublicKey{}

	for i, v := range b.nextEpochShardCandidate {
		key, err := v.ToBase58()
		if err != nil {
			return newCommitteeChange, err
		}
		indexNextEpochShardCandidate[key] = i
	}

	for _, committeePublicKey := range unstakeInstruction.PublicKeys {
		if common.IndexOfStr(committeePublicKey, env.subtituteCandidates) == -1 {
			if common.IndexOfStr(committeePublicKey, env.validators) == -1 {

				// if not found then delete auto staking data for this public key if present
				if _, ok := b.autoStake[committeePublicKey]; ok {
					delete(b.autoStake, committeePublicKey)
				}

			} else {
				Logger.log.Info("[unstake] Node is in active consensus state")
				// if found in committee list then turn off auto staking
				if _, ok := b.autoStake[committeePublicKey]; ok {
					b.autoStake[committeePublicKey] = false
					committeeChange.Unstake = append(committeeChange.Unstake, committeePublicKey)
				}
			}

		} else {
			Logger.log.Info("[unstake] Node is in next epoch candidate")

			if _, ok := b.autoStake[committeePublicKey]; ok {
				delete(b.autoStake, committeePublicKey)
			}

			if _, ok := b.stakingTx[committeePublicKey]; ok {
				delete(b.stakingTx, committeePublicKey)
			}

			committeePublicKeyStruct := incognitokey.CommitteePublicKey{}
			err := committeePublicKeyStruct.FromBase58(committeePublicKey)
			if err != nil {
				return committeeChange, err
			}

			rewardReceiverKey := committeePublicKeyStruct.GetIncKeyBase58()
			if _, ok := b.rewardReceiver[rewardReceiverKey]; ok {
				delete(b.rewardReceiver, rewardReceiverKey)
			}

			//TODO: Scale by checking unstake instruction type later
			newCommitteeChange.NextEpochShardCandidateRemoved =
				append(newCommitteeChange.NextEpochShardCandidateRemoved, committeePublicKeyStruct)

			indexCandidate := indexNextEpochShardCandidate[committeePublicKey]
			b.nextEpochShardCandidate = append(b.nextEpochShardCandidate[:indexCandidate], b.nextEpochShardCandidate[indexCandidate+1:]...)
			nextEpochShardCandidateUnstakeKey = append(nextEpochShardCandidateUnstakeKey, committeePublicKeyStruct)
		}
	}

	err := statedb.DeleteStakerInfo(env.ConsensusStateDB, nextEpochShardCandidateUnstakeKey)
	return newCommitteeChange, err
}

func (b *BeaconCommitteeStateV1) getSubtituteCandidates() ([]string, error) {
	commonPoolValidators := []string{}

	candidateBeaconWaitingForNextRandom := b.nextEpochBeaconCandidate
	candidateBeaconWaitingForNextRandomStr, err := incognitokey.CommitteeKeyListToString(candidateBeaconWaitingForNextRandom)
	if err != nil {
		return nil, err
	}
	commonPoolValidators = append(commonPoolValidators, candidateBeaconWaitingForNextRandomStr...)
	candidateShardWaitingForNextRandom := b.nextEpochShardCandidate
	candidateShardWaitingForNextRandomStr, err := incognitokey.CommitteeKeyListToString(candidateShardWaitingForNextRandom)
	if err != nil {
		return nil, err
	}
	commonPoolValidators = append(commonPoolValidators, candidateShardWaitingForNextRandomStr...)

	return commonPoolValidators, nil
}

func (b *BeaconCommitteeStateV1) getValidators() ([]string, error) {
	validators := []string{}

	for _, committee := range b.shardCommittee {
		committeeStr, err := incognitokey.CommitteeKeyListToString(committee)
		if err != nil {
			return nil, err
		}
		validators = append(validators, committeeStr...)
	}
	for _, substitute := range b.shardSubstitute {
		substituteStr, err := incognitokey.CommitteeKeyListToString(substitute)
		if err != nil {
			return nil, err
		}
		validators = append(validators, substituteStr...)
	}

	beaconCommittee := b.beaconCommittee
	beaconCommitteeStr, err := incognitokey.CommitteeKeyListToString(beaconCommittee)
	if err != nil {
		return nil, err
	}
	validators = append(validators, beaconCommitteeStr...)
	beaconSubstitute := b.beaconSubstitute
	beaconSubstituteStr, err := incognitokey.CommitteeKeyListToString(beaconSubstitute)
	if err != nil {
		return nil, err
	}
	validators = append(validators, beaconSubstituteStr...)

	candidateBeaconWaitingForCurrentRandom := b.currentEpochBeaconCandidate
	candidateBeaconWaitingForCurrentRandomStr, err := incognitokey.CommitteeKeyListToString(candidateBeaconWaitingForCurrentRandom)
	if err != nil {
		return nil, err
	}
	validators = append(validators, candidateBeaconWaitingForCurrentRandomStr...)
	candidateShardWaitingForCurrentRandom := b.currentEpochShardCandidate
	candidateShardWaitingForCurrentRandomStr, err := incognitokey.CommitteeKeyListToString(candidateShardWaitingForCurrentRandom)
	if err != nil {
		return nil, err
	}
	validators = append(validators, candidateShardWaitingForCurrentRandomStr...)

	return validators, nil
}

func (b *BeaconCommitteeStateV1) processUnstakeChange(committeeChange *CommitteeChange, env *BeaconCommitteeStateEnvironment) (*CommitteeChange, error) {

	newCommitteeChange := committeeChange

	unstakingIncognitoKey, err := incognitokey.CommitteeBase58KeyListToStruct(committeeChange.Unstake)
	if err != nil {
		return newCommitteeChange, err
	}
	err = statedb.StoreStakerInfo(
		env.ConsensusStateDB,
		unstakingIncognitoKey,
		b.rewardReceiver,
		b.autoStake,
		b.stakingTx,
	)
	return newCommitteeChange, err
}
