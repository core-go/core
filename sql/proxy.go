package sql

import "context"

type Proxy interface {
	BeginTransaction(ctx context.Context, timeout int64) (string, error)
	CommitTransaction(ctx context.Context, tx string) error
	RollbackTransaction(ctx context.Context, tx string) error
	Exec(ctx context.Context, query string, values ...interface{}) (int64, error)
	ExecBatch(ctx context.Context, master bool, stm ...Statement) (int64, error)
	Query(ctx context.Context, result interface{}, query string, values ...interface{}) error
	QueryOne(ctx context.Context, result interface{}, query string, values ...interface{}) error
	ExecTx(ctx context.Context, tx string, commit bool, query string, values ...interface{}) (int64, error)
	ExecBatchTx(ctx context.Context, tx string, commit bool, master bool, stm ...Statement) (int64, error)
	QueryTx(ctx context.Context, tx string, commit bool, result interface{}, query string, values ...interface{}) error
	QueryOneTx(ctx context.Context, tx string, commit bool, result interface{}, query string, values ...interface{}) error

	Insert(ctx context.Context, table string, model interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	Update(ctx context.Context, table string, model interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	Save(ctx context.Context, table string, model interface{}, driver string, options...*Schema) (int64, error)
	InsertBatch(ctx context.Context, table string, models interface{}, driver string, options...*Schema) (int64, error)
	UpdateBatch(ctx context.Context, table string, models interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	SaveBatch(ctx context.Context, table string, models interface{}, driver string, options...*Schema) (int64, error)

	InsertTx(ctx context.Context, tx string, commit bool, table string, model interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	UpdateTx(ctx context.Context, tx string, commit bool, table string, model interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	SaveTx(ctx context.Context, tx string, commit bool, table string, model interface{}, driver string, options...*Schema) (int64, error)
	InsertBatchTx(ctx context.Context, tx string, commit bool, table string, models interface{}, driver string, options...*Schema) (int64, error)
	UpdateBatchTx(ctx context.Context, tx string, commit bool, table string, models interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	SaveBatchTx(ctx context.Context, tx string, commit bool, table string, models interface{}, driver string, options...*Schema) (int64, error)

	InsertAndCommit(ctx context.Context, tx string, table string, model interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	UpdateAndCommit(ctx context.Context, tx string, table string, model interface{}, driver string, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	SaveAndCommit(ctx context.Context, tx string, table string, model interface{}, driver string, options...*Schema) (int64, error)
	InsertBatchAndCommit(ctx context.Context, tx string, table string, models interface{}, driver string, options...*Schema) (int64, error)
	UpdateBatchAndCommit(ctx context.Context, tx string, table string, models interface{}, buildParam func(int) string, boolSupport bool, options...*Schema) (int64, error)
	SaveBatchAndCommit(ctx context.Context, tx string, table string, models interface{}, driver string, options...*Schema) (int64, error)
}
const Timeout int64 = 5000000000
func BeginTx(ctx context.Context, proxy Proxy, timeouts... int64) (context.Context, string, error) {
	timeout := Timeout
	if len(timeouts) > 0 && timeouts[0] > 0 {
		timeout = timeouts[0]
	}
	tx, err := proxy.BeginTransaction(ctx, timeout)
	if err != nil {
		return ctx, tx, err
	}
	c2 := context.WithValue(ctx, "txId", &tx)
	return c2, tx, nil
}
func CommitTx(ctx context.Context, proxy Proxy, tx string, err error, options...bool) error {
	if err != nil {
		if !(len(options) > 0 && options[0] == false) {
			er := proxy.RollbackTransaction(ctx, tx)
			if er != nil {
				return er
			}
		}
		return err
	}
	return proxy.CommitTransaction(ctx, tx)
}
func EndTx(ctx context.Context, proxy Proxy, tx string, res int64, err error, options...bool) (int64, error) {
	er := CommitTx(ctx, proxy, tx, err, options...)
	return res, er
}
func ExecProxy(ctx context.Context, proxy Proxy, query string, args ...interface{}) (int64, error) {
	tx := GetTxId(ctx)
	if tx == nil {
		return proxy.Exec(ctx, query, args...)
	}
	return proxy.ExecTx(ctx, *tx, false, query, args...)
}
func QueryProxy(ctx context.Context, proxy Proxy, result interface{}, query string, args ...interface{}) error {
	tx := GetTxId(ctx)
	if tx == nil {
		return proxy.Query(ctx, result, query, args...)
	}
	return proxy.QueryTx(ctx, *tx, false, result, query, args...)
}
func QueryOneProxy(ctx context.Context, proxy Proxy, result interface{}, query string, args ...interface{}) error {
	tx := GetTxId(ctx)
	if tx == nil {
		return proxy.QueryOne(ctx, result, query, args...)
	}
	return proxy.QueryOneTx(ctx, *tx, false, result, query, args...)
}
