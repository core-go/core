package settings

import (
	"context"
	"database/sql"
	"fmt"
)

type UseCase struct {
	DB *sql.DB
	BuildParam func(int)string
	Table string
	Id  string
	Language string
	DateFormat string
}

func NewSettingsService(db *sql.DB, buildParam func(int)string, table string, opts...string) *UseCase {
	var id, language, dateFormat string
	if len(opts) > 0 {
		id = opts[0]
	} else {
		id = "id"
	}
	if len(opts) > 1 {
		dateFormat = opts[1]
	} else {
		dateFormat = "dateformat"
	}
	if len(opts) > 2 {
		language = opts[2]
	} else {
		language = "language"
	}
	return &UseCase{DB: db, BuildParam: buildParam, Table: table, Id: id, DateFormat: dateFormat, Language: language}
}

func (a *UseCase) Save(ctx context.Context, id string, settings Settings) (int64, error) {
	query := fmt.Sprintf("update %s set %s = %s, %s = %s where %s = %s",
		a.Table, a.Language, a.BuildParam(1), a.DateFormat, a.BuildParam(2), a.Id, a.BuildParam(3))
	res, err := a.DB.ExecContext(ctx, query, settings.Language, settings.DateFormat, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
