package qdb

import (
	"fmt"

	protos "github.com/pg-sharding/spqr/pkg/protos"
)

const (
	CMD_PUT    string = "put"
	CMD_DELETE string = "delete"
)

type QdbStatement struct {
	CmdType   string
	Key       string
	Value     string
	Extension string
}

func NewQdbStatement(cmdType string, key string, value string) (*QdbStatement, error) {
	if cmdType != CMD_PUT && cmdType != CMD_DELETE {
		return nil, fmt.Errorf("unknown type of QdbStatement: %s", cmdType)
	}
	return &QdbStatement{CmdType: cmdType, Key: key, Value: value}, nil
}

func NewQdbStatementExt(cmdType string, key string, value string, extension string) (*QdbStatement, error) {
	if stmt, err := NewQdbStatement(cmdType, key, value); err != nil {
		return nil, err
	} else {
		stmt.Extension = extension
		return stmt, nil
	}
}

func (s *QdbStatement) ToProto() *protos.QdbTransactionCmd {
	return &protos.QdbTransactionCmd{Command: s.CmdType, Key: s.Key, Value: s.Value}
}

func SliceToProto(stmts []QdbStatement) []*protos.QdbTransactionCmd {
	result := make([]*protos.QdbTransactionCmd, len(stmts))
	for i, qdbStmt := range stmts {
		result[i] = qdbStmt.ToProto()
	}
	return result
}

func QdbStmtFromProto(protoStmt *protos.QdbTransactionCmd) (*QdbStatement, error) {
	return NewQdbStatement(protoStmt.Command, protoStmt.Key, protoStmt.Value)
}

type QdbTransaction struct {
	transactionId string
	commands      []QdbStatement
}

func (t *QdbTransaction) Id() string {
	return t.transactionId
}

func NewTransaction(transactionId string) (*QdbTransaction, error) {
	if len(transactionId) == 0 {
		return nil, fmt.Errorf("infalid transaction id")
	}
	return &QdbTransaction{transactionId: transactionId, commands: make([]QdbStatement, 0)}, nil
}
func (t *QdbTransaction) Append(qdbExecutors []QdbStatement) error {
	if len(qdbExecutors) == 0 {
		return fmt.Errorf("cant't add empty list of DB changes to transaction %s", t.transactionId)
	}
	t.commands = append(t.commands, qdbExecutors...)
	return nil
}

func (t *QdbTransaction) Validate() error {
	if len(t.commands) == 0 {
		return fmt.Errorf("transaction %s haven't statements", t.transactionId)
	}
	return nil
}
