package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct {
	Client *sql.DB
}

func (mc *MysqlClient) increment(table string, valueColumn string) (int64, error) {
	stmt := fmt.Sprintf(
		"UPDATE %s SET %s = LAST_INSERT_ID(%s + ?)",
		table,
		valueColumn,
		valueColumn,
	)

	res, errExec := mc.Client.Exec(stmt)
	if errExec != nil {
		return 0, errExec
	}
	incr, err := res.LastInsertId()
	return incr, err
}

func (mc *MysqlClient) close() {
	mc.Client.Close()
}

// MysqlNewClient クライアントを返す
func MysqlNewClient(host string, user string, password string, port string, dbName string) MysqlClient {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&autocommit=true",
		user,
		password,
		host,
		port,
		dbName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return MysqlClient{Client: db}
}
