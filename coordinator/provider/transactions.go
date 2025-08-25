package provider

import (
	"context"

	"github.com/pg-sharding/spqr/coordinator"
	mtran "github.com/pg-sharding/spqr/pkg/models/transaction"
	protos "github.com/pg-sharding/spqr/pkg/protos"
	qdb "github.com/pg-sharding/spqr/qdb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetaTransactionServer struct {
	protos.UnimplementedMetaTransactionServiceServer

	impl coordinator.Coordinator
}

func NewMetaTransactionServer(impl coordinator.Coordinator) *MetaTransactionServer {
	return &MetaTransactionServer{
		impl: impl,
	}
}

var _ protos.MetaTransactionServiceServer = &MetaTransactionServer{}

func toQdbStatementList(cmdList []*protos.QdbTransactionCmd) ([]qdb.QdbStatement, error) {
	stmts := make([]qdb.QdbStatement, 0, len(cmdList))
	for _, cmd := range cmdList {
		if qdbCmd, err := qdb.QdbStmtFromProto(cmd); err != nil {
			return nil, err
		} else {
			stmts = append(stmts, *qdbCmd)
		}
	}
	return stmts, nil
}

func (mts *MetaTransactionServer) ExecNoTran(ctx context.Context, request *protos.ExecNoTranRequest) (*emptypb.Empty, error) {
	if stmts, err := toQdbStatementList(request.CmdList); err != nil {
		return nil, err
	} else {
		if tranChunk, err := mtran.NewMetaTransactionChunk(request.MetaCmdList, stmts); err != nil {
			return nil, err
		} else {
			return nil, mts.impl.ExecNoTran(ctx, tranChunk)
		}
	}
}
func (mts *MetaTransactionServer) ExecTran(ctx context.Context, request *protos.MetaTransactionRequest) (*emptypb.Empty, error) {
	if metaTran, err := mtran.NewTransaction(request.Transaction.TransactionId); err != nil {
		return nil, err
	} else {
		if stmts, err := toQdbStatementList(request.Transaction.CmdList); err != nil {
			return nil, err
		} else {
			if err := metaTran.AppendStatements(request.Transaction.MetaCmdList, stmts); err != nil {
				return nil, err
			}
		}
		return nil, mts.impl.ExecTran(ctx, metaTran)
	}
}
