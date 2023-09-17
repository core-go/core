package sql

import (
	"context"
	"database/sql"
	"fmt"
)

type PrivilegesLoader struct {
	DB    *sql.DB
	Query string
}

func NewPrivilegesLoader(db *sql.DB, query string, options ...bool) *PrivilegesLoader {
	var handleDriver bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	if handleDriver {
		driver := getDriver(db)
		query = replaceQueryArgs(driver, query)
	}
	return &PrivilegesLoader{DB: db, Query: query}
}

func (l PrivilegesLoader) Privileges(ctx context.Context, userId string) []string {
	privileges := make([]string, 0)
	rows, err := l.DB.QueryContext(ctx, l.Query, userId)
	if err != nil {
		return privileges
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var permissions int32
		if err = rows.Scan(&id, &permissions); err == nil {
			if permissions != actionNone {
				x := id + " " + fmt.Sprintf("%X", permissions)
				privileges = append(privileges, x)
			} else {
				privileges = append(privileges, id)
			}
		}
	}
	return privileges
}
