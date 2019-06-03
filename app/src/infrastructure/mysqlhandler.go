package infrastructure

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"interfaces"
	"os"
	"time"
)

var schemaCounty = `
create table if not exists country
(
    id varchar(2) not null,
    name varchar(255) not null,
    CONSTRAINT country_pk PRIMARY KEY (id)
)
comment 'Name of countries' DEFAULT CHARSET=utf8 DEFAULT COLLATE utf8_unicode_ci;
`

var schemaCode = `
create table if not exists code
(
    id varchar(2) not null,
    name varchar(255) not null,
    CONSTRAINT code_pk PRIMARY KEY (id)
)
comment 'Phone codes of country' DEFAULT CHARSET=utf8 DEFAULT COLLATE utf8_unicode_ci;
`

type MysqlHandler struct {
	Conn   *sqlx.DB
	logger *Logger
}

func (handler *MysqlHandler) Execute(statement string, model interface{}) {
	handler.Conn.NamedExec(statement, model)
}

func (handler *MysqlHandler) Query(statement string, params map[string]interface{}) interfaces.Row {
	rows, err := handler.Conn.NamedQuery(statement, params)
	if err != nil {
		fmt.Println(err)
		return new(MysqlRow)
	}
	row := new(MysqlRow)
	row.Rows = rows
	return row
}

func (handler *MysqlHandler) Trunc(table string) {
	handler.Conn.MustExec("TRUNCATE " + table)
}

type MysqlRow struct {
	Rows *sqlx.Rows
}

func (r MysqlRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

func (r MysqlRow) Next() bool {
	return r.Rows.Next()
}

func NewMysqlHandler(connect string, logger *Logger) *MysqlHandler {
	var db *sqlx.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = sqlx.Connect("mysql", connect)
		if err == nil {
			break
		} else {
			logger.Error("Ошибка подключения " + err.Error())
		}

		time.Sleep(1 * time.Second)
	}

	if db == nil {
		os.Exit(1)
	}

	db.MustExec(schemaCounty)
	db.MustExec(schemaCode)

	mysqlHandler := new(MysqlHandler)
	mysqlHandler.Conn = db
	mysqlHandler.logger = logger
	return mysqlHandler
}
