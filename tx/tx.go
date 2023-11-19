package tx

import (
	"context"
	"database/sql"
	"errors"
	"runtime/debug"
)
func Callback(ctx context.Context, db *sql.DB, callback func(context.Context)error, opts ...string) (err error) {
	txName := "tx"
	if len(opts) > 0 {
		txName = opts[0]
	}
	return CallbackTx(ctx, db, txName, callback)
}
func CallbackTx(ctx context.Context, db *sql.DB, txName string, callback func(context.Context)error, opts ...*sql.TxOptions) (err error) {
	var tx *sql.Tx
	if len(opts) > 0 && opts[0] != nil {
		tx, err = db.BeginTx(ctx, opts[0])
	} else {
		tx, err = db.BeginTx(ctx, nil)
	}
	if err != nil {
		return err
	}
	defer func(e error) {
		if err0 := recover(); err0 != nil {
			tx.Rollback()
			debug.PrintStack()
			err = errors.New("error when execute sql")
			return
		}
	}(err)
	ctx = context.WithValue(ctx, txName, tx)
	if err = callback(ctx); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func Execute(ctx context.Context, db *sql.DB, callback func(context.Context)(int64, error), opts ...string) (int64, error) {
	txName := "tx"
	if len(opts) > 0 {
		txName = opts[0]
	}
	return ExecuteTx(ctx, db, txName, callback)
}
func ExecuteTx(ctx context.Context, db *sql.DB, txName string, callback func(context.Context)(int64, error), opts ...*sql.TxOptions) (int64, error) {
	var res int64
	er0 := CallbackTx(ctx, db, txName, func(context.Context) error {
		result, err := callback(ctx)
		if err != nil {
			return err
		}
		res = result
		return nil
	}, opts...)
	return res, er0
}
