package evmtask

import (
	"context"

	"github.com/Spacescore/observatory-task/pkg/errors"
	"github.com/Spacescore/observatory-task/pkg/lotus"
	"github.com/Spacescore/observatory-task/pkg/models/evmmodel"
	"github.com/Spacescore/observatory-task/pkg/storage"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	log "github.com/sirupsen/logrus"
)

// BlockHeader extract block header for evm
type BlockHeader struct {
}

func (b *BlockHeader) Name() string {
	return "evm_block_header"
}

func (b *BlockHeader) Model() interface{} {
	return new(evmmodel.BlockHeader)
}

func (b *BlockHeader) Run(ctx context.Context, rpc *lotus.Rpc, version int, tipSet *types.TipSet, force bool, storage storage.Storage) error {
	parentTs, err := rpc.Node().ChainGetTipSet(ctx, tipSet.Parents())
	if err != nil {
		return errors.Wrap(err, "ChainGetTipSet failed")
	}

	if !force {
		existed, err := storage.Existed(b.Model(), int64(parentTs.Height()), version)
		if err != nil {
			return errors.Wrap(err, "storage.Existed failed")
		}
		if existed {
			log.Infof("task [%s] has been process (%d,%d), ignore it", b.Name(), int64(parentTs.Height()), version)
			return nil
		}
	}

	tipSetCid, err := parentTs.Key().Cid()
	if err != nil {
		return errors.Wrap(err, "tipSetCid failed")
	}

	hash, err := ethtypes.EthHashFromCid(tipSetCid)
	if err != nil {
		return errors.Wrap(err, "rpc EthHashFromCid failed")
	}

	ethBlock, err := rpc.Node().EthGetBlockByHash(ctx, hash, false)
	if err != nil {
		return errors.Wrap(err, "rpc EthGetBlockByHash failed")
	}
	if ethBlock.Number == 0 {
		log.Infof("block number == 0")
		return nil
	}

	blockHeader := &evmmodel.BlockHeader{
		Height:           int64(parentTs.Height()),
		Version:          version,
		Hash:             hash.String(),
		ParentHash:       ethBlock.ParentHash.String(),
		Miner:            ethBlock.Miner.String(),
		StateRoot:        ethBlock.StateRoot.String(),
		TransactionsRoot: ethBlock.TransactionsRoot.String(),
		ReceiptsRoot:     ethBlock.ReceiptsRoot.String(),
		Difficulty:       int64(ethBlock.Difficulty),
		Number:           int64(ethBlock.Number),
		GasLimit:         int64(ethBlock.GasLimit),
		GasUsed:          int64(ethBlock.GasUsed),
		Timestamp:        int64(ethBlock.Timestamp),
		ExtraData:        string(ethBlock.Extradata),
		MixHash:          ethBlock.MixHash.String(),
		Nonce:            ethBlock.Nonce.String(),
		BaseFeePerGas:    ethBlock.BaseFeePerGas.String(),
		Size:             uint64(ethBlock.Size),
		Sha3Uncles:       ethBlock.Sha3Uncles.String(),
	}

	if err = storage.Insert(ctx, new(evmmodel.BlockHeader),
		int64(parentTs.Height()), version, blockHeader); err != nil {
		return errors.Wrap(err, "storageWrite failed")
	}

	log.Infof("Tipset[%v] has been process height: %v evm_block_header", tipSet.Height(), tipSet.Height())
	return nil
}
