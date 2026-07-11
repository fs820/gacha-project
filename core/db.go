package core

import (
	"database/sql" // データベース操作
	"log"          // ロギング
	"os"           // 環境変数の取得

	_ "github.com/lib/pq" // PostgreSQLドライバ (データベース接続)
)

// DB
var UserDB *sql.DB

// データベースの初期化関数
func InitDB() {
	var err error
	// .env や Render の環境変数からURLを取得
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	// PostgreSQLに接続
	UserDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// ユーザーのカウンター状態を保存するテーブルを作成
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		uid TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password_hash TEXT,
		stones INTEGER DEFAULT 30000,
		star4_limit_counter INTEGER DEFAULT 0,
		star5_limit_counter INTEGER DEFAULT 0,
		is_next_pickup_guaranteed BOOLEAN DEFAULT false
	);`
	_, err = UserDB.Exec(usersTable)
	if err != nil {
		log.Fatal("usersテーブル作成エラー:", err)
	}

	// ガチャの履歴を保存するテーブルを作成
	historyTable := `
	CREATE TABLE IF NOT EXISTS history (
		id SERIAL PRIMARY KEY,
		uid TEXT,
		rarity TEXT,
		character TEXT
	);`
	_, err = UserDB.Exec(historyTable)
	if err != nil {
		log.Fatal("historyテーブル作成エラー:", err)
	}

	// キャラクターデータを保存するテーブルを作成
	charactersTable := `
	CREATE TABLE IF NOT EXISTS characters (
		id SERIAL PRIMARY KEY,
		name TEXT,
		rarity TEXT,
		is_pickup BOOLEAN DEFAULT false
	);`
	_, err = UserDB.Exec(charactersTable)
	if err != nil {
		log.Fatal("charactersテーブル作成エラー:", err)
	}

	// もしキャラクターテーブルが空の場合初期化
	var count int
	UserDB.QueryRow("SELECT COUNT(*) FROM characters").Scan(&count)
	if count == 0 {
		log.Println("キャラクターの初期データを挿入します...")
		initialData := []struct {
			name     string
			rarity   string
			isPickup bool
		}{
			{"ゼウス", "星5", true},
			{"ウラノス", "星5", false},
			{"クロノス", "星5", false},
			{"釈迦", "星5", false},
			{"シヴァ", "星5", false},
			{"ポセイドン", "星5", false},
			{"ヘラクレス", "星5", false},
			{"キリスト", "星5", false},
			{"ヨハネ", "星4", true},
			{"千手観音", "星4", true},
			{"アキレス", "星4", true},
			{"武器", "星3", false},
		}

		for _, c := range initialData {
			UserDB.Exec("INSERT INTO characters (name, rarity, is_pickup) VALUES ($1, $2, $3)",
				c.name, c.rarity, c.isPickup)
		}
	}

	// 決済の注文状態を管理するテーブルを作成
	ordersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		order_id TEXT PRIMARY KEY,
		uid TEXT,
		amount INTEGER,
		status TEXT   /* 'pending'(未払い) または 'paid'(支払い済み) */
	);`
	_, err = UserDB.Exec(ordersTable)
	if err != nil {
		log.Fatal("ordersテーブル作成エラー:", err)
	}
}
