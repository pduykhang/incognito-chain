package blsbft

import (
	"reflect"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/types"
	"github.com/incognitochain/incognito-chain/common"
	signatureschemes2 "github.com/incognitochain/incognito-chain/consensus_v2/signatureschemes"
	"github.com/incognitochain/incognito-chain/incognitokey"
	"github.com/incognitochain/incognito-chain/multiview"
)

func Test_actorV2_handleProposeMsg(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		proposeMsg BFTPropose
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				proposeMsg: BFTPropose{
					PeerID:   "",
					Block:    nil,
					TimeSlot: 19,
				},
			},
			wantErr: true,
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				proposeMsg: BFTPropose{
					PeerID:   "",
					Block:    nil,
					TimeSlot: 19,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			if err := actorV2.handleProposeMsg(tt.args.proposeMsg); (err != nil) != tt.wantErr {
				t.Errorf("actorV2.handleProposeMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_actorV2_handleVoteMsg(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		voteMsg BFTVote
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			if err := actorV2.handleVoteMsg(tt.args.voteMsg); (err != nil) != tt.wantErr {
				t.Errorf("actorV2.handleVoteMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_actorV2_proposeBlock(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		userMiningKey     signatureschemes2.MiningKey
		proposerPk        incognitokey.CommitteePublicKey
		block             types.BlockInterface
		committees        []incognitokey.CommitteePublicKey
		committeeViewHash common.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.BlockInterface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			got, err := actorV2.proposeBlock(tt.args.userMiningKey, tt.args.proposerPk, tt.args.block, tt.args.committees, tt.args.committeeViewHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("actorV2.proposeBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("actorV2.proposeBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_actorV2_proposeBeaconBlock(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		b58Str            string
		block             types.BlockInterface
		committees        []incognitokey.CommitteePublicKey
		committeeViewHash common.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.BlockInterface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			got, err := actorV2.proposeBeaconBlock(tt.args.b58Str, tt.args.block, tt.args.committees, tt.args.committeeViewHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("actorV2.proposeBeaconBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("actorV2.proposeBeaconBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_actorV2_proposeShardBlock(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		b58Str            string
		block             types.BlockInterface
		committees        []incognitokey.CommitteePublicKey
		committeeViewHash common.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.BlockInterface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			got, err := actorV2.proposeShardBlock(tt.args.b58Str, tt.args.block, tt.args.committees, tt.args.committeeViewHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("actorV2.proposeShardBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("actorV2.proposeShardBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_actorV2_getValidProposeBlocks(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		bestView multiview.View
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*ProposeBlockInfo
	}{
		{
			name: "",
			fields: fields{
				actorBase:            actorBase{},
				committeeChain:       &blockchain.BeaconChain{},
				currentTime:          1,
				currentTimeSlot:      1,
				proposeHistory:       &lru.Cache{},
				receiveBlockByHeight: map[uint64][]*ProposeBlockInfo{},
				receiveBlockByHash:   map[string]*ProposeBlockInfo{},
				blockVersion:         1,
			},
			args: args{
				bestView: &blockchain.BeaconBestState{},
			},
			want: []*ProposeBlockInfo{
				&ProposeBlockInfo{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			if got := actorV2.getValidProposeBlocks(tt.args.bestView); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("actorV2.getValidProposeBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_actorV2_validateBlock(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		bestView         multiview.View
		proposeBlockInfo *ProposeBlockInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			if err := actorV2.validateBlock(tt.args.bestView, tt.args.proposeBlockInfo); (err != nil) != tt.wantErr {
				t.Errorf("actorV2.validateBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_actorV2_voteForBlock(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		v *ProposeBlockInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			if err := actorV2.voteForBlock(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("actorV2.voteForBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_actorV2_processIfBlockGetEnoughVote(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		blockHash string
		v         *ProposeBlockInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			actorV2.processIfBlockGetEnoughVote(tt.args.blockHash, tt.args.v)
		})
	}
}

func Test_actorV2_processWithEnoughVotes(t *testing.T) {
	type fields struct {
		actorBase            actorBase
		committeeChain       blockchain.Chain
		currentTime          int64
		currentTimeSlot      int64
		proposeHistory       *lru.Cache
		receiveBlockByHeight map[uint64][]*ProposeBlockInfo
		receiveBlockByHash   map[string]*ProposeBlockInfo
		voteHistory          map[uint64]types.BlockInterface
		bodyHashes           map[uint64]map[string]bool
		votedTimeslot        map[int64]bool
		blockVersion         int
	}
	type args struct {
		v *ProposeBlockInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorV2 := &actorV2{
				actorBase:            tt.fields.actorBase,
				committeeChain:       tt.fields.committeeChain,
				currentTime:          tt.fields.currentTime,
				currentTimeSlot:      tt.fields.currentTimeSlot,
				proposeHistory:       tt.fields.proposeHistory,
				receiveBlockByHeight: tt.fields.receiveBlockByHeight,
				receiveBlockByHash:   tt.fields.receiveBlockByHash,
				voteHistory:          tt.fields.voteHistory,
				bodyHashes:           tt.fields.bodyHashes,
				votedTimeslot:        tt.fields.votedTimeslot,
				blockVersion:         tt.fields.blockVersion,
			}
			if err := actorV2.processWithEnoughVotes(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("actorV2.processWithEnoughVotes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
