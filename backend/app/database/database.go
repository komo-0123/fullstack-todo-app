package database

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func Init() error {
	var err error
	// dsn -> ユーザー名:パスワード@tcp(ホスト名:ポート番号)/データベース名?オプション
	dsn := "todo_user:todo_password@tcp(mysql-container:3306)/todo_db"
	db, err = sql.Open("mysql", dsn) // DBとの接続を準備
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	// DBへの接続を確認
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
