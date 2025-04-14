package db

import (
	"FGW/internal/config"
	"FGW/pkg/wlogger/msg"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"log"
)

// MSSQLConn подключение к БД.
func MSSQLConn(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("Driver=%s; Server=%s,%s; Database=%s; Uid=%s; Pwd=%s; ClientCharset=WINDOWS-1251",
		cfg.MSSQL.MSSQLDriver,
		cfg.MSSQL.MSSQLServer,
		cfg.MSSQL.MSSQLPort,
		cfg.MSSQL.MSSQLDatabase,
		cfg.MSSQL.MSSQLUsername,
		cfg.MSSQL.MSSQLPassword,
	)

	mssqlConn, err := sql.Open("odbc", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = mssqlConn.PingContext(ctx); err != nil {
		CloseDB(mssqlConn)
		return nil, err
	}

	log.Printf("%s", msg.I2000)
	return mssqlConn, nil
}

// CloseDB закрывает соединение с БД после выполнения транзакции.
func CloseDB(dbConn *sql.DB) {
	if err := dbConn.Close(); err != nil {
		return
	}
	log.Printf("%s", msg.I2001)
}

// CloseRows закрывает строки с данными.
func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		return
	}
	log.Printf("%s", msg.I2002)
}
