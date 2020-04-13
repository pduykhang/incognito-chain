package rpcservice

import (
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/syncker"
)

type SynkerService struct {
	Synker *syncker.SynckerManager
}

func (s *SynkerService) GetBeaconPoolInfo() []common.BlockPoolInterface {
	return s.Synker.GetPoolInfo(syncker.BeaconPoolType, 0)
}

func (s *SynkerService) GetShardPoolInfo(shardID int) []common.BlockPoolInterface {
	return s.Synker.GetPoolInfo(syncker.ShardPoolType, shardID)
}

func (s *SynkerService) GetCrossShardPoolInfo(toShard int) []common.BlockPoolInterface {
	return s.Synker.GetPoolInfo(syncker.CrossShardPoolType, toShard)
}

func (s *SynkerService) GetShardToBeaconPoolInfo() []common.BlockPoolInterface {
	return s.Synker.GetPoolInfo(syncker.S2BPoolType, 0)
}
