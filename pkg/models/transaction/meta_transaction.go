package meta_transaction

import (
	"context"
	"fmt"

	proto "github.com/pg-sharding/spqr/pkg/protos"
	"github.com/pg-sharding/spqr/qdb"
)

type TransactionExecutor interface {
	ExecNoTran(ctx context.Context, chunk *MetaTransactionChunk) error
	ExecTran(ctx context.Context, transaction *MetaTransaction) error
}

type MetaTransaction struct {
	GossipRequests []*proto.MetaTransactionGossip
	QdbTransaction qdb.QdbTransaction
}

type MetaTransactionChunk struct {
	GossipRequests []*proto.MetaTransactionGossip
	QdbStatements  []qdb.QdbStatement
}

const (
	GR_UNKNOWN = iota
	GR_CreateDistributionRequest
	GR_CreateKeyRange
)

func NewMetaTransactionChunk(gossipRequests []*proto.MetaTransactionGossip,
	qdbStatements []qdb.QdbStatement) (*MetaTransactionChunk, error) {
	if len(qdbStatements) == 0 {
		return nil, fmt.Errorf("transaction chunk must have a qdb statetment (case 0)")
	} else {
		return &MetaTransactionChunk{
			GossipRequests: gossipRequests,
			QdbStatements:  qdbStatements,
		}, nil
	}
}

func (tc *MetaTransactionChunk) Append(gossipRequests []*proto.MetaTransactionGossip,
	qdbStatements []qdb.QdbStatement) error {
	if len(qdbStatements) == 0 {
		return fmt.Errorf("transaction chunk must have a qdb statetment (case 1)")
	} else {
		tc.GossipRequests = append(tc.GossipRequests, gossipRequests...)
		tc.QdbStatements = append(tc.QdbStatements, qdbStatements...)
	}
	return nil
}

func NewTransaction(transactionId string) (*MetaTransaction, error) {
	if qdbTran, err := qdb.NewTransaction(transactionId); err != nil {
		return nil, err
	} else {
		return &MetaTransaction{
			GossipRequests: make([]*proto.MetaTransactionGossip, 0),
			QdbTransaction: *qdbTran,
		}, nil
	}
}

func GetGossipRequestType(request *proto.MetaTransactionGossip) int {
	result := GR_UNKNOWN
	if request.CreateDistribution != nil {
		result = GR_CreateDistributionRequest
	}
	if request.CreateKeyRangeRequest != nil {
		if result != GR_UNKNOWN {
			return GR_UNKNOWN
		} else {
			result = GR_CreateKeyRange
		}
	}
	return result
}

func (mt *MetaTransaction) AppendStatements(gossipRequests []*proto.MetaTransactionGossip,
	qdbStatements []qdb.QdbStatement) error {
	if err := mt.QdbTransaction.Append(qdbStatements); err != nil {
		return err
	} else {
		for _, req := range gossipRequests {
			if gossipType := GetGossipRequestType(req); gossipType == GR_UNKNOWN {
				return fmt.Errorf("invalid meta transaction request (case 0)")
			}
			mt.GossipRequests = append(mt.GossipRequests, req)
		}
		return nil
	}
}
