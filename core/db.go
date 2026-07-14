package core

import (
	"database/sql" // データベース操作
	"log"          // ロギング
	"os"           // 環境変数の取得

	_ "github.com/lib/pq" // PostgreSQLドライバ (データベース接続)
)

// 環境変数からDB接続を生成する関数
func NewDBConnection() (*sql.DB, error) {
	// .env や Render の環境変数からURLを取得
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	// PostgreSQLに接続
	return sql.Open("postgres", dbURL)
}

// データベース初期化関数
func InitSchema(db *sql.DB) error {
	var err error

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
	if _, err = db.Exec(usersTable); err != nil {
		return err
	}

	// キャラクターデータを保存するテーブルを作成
	charactersTable := `
	CREATE TABLE IF NOT EXISTS characters (
		id SERIAL PRIMARY KEY,
        name TEXT UNIQUE NOT NULL,
		rarity TEXT NOT NULL
	);`
	if _, err = db.Exec(charactersTable); err != nil {
		return err
	}

	// もしキャラクターテーブルが空の場合恒常キャラクターを挿入
	var count int
	db.QueryRow("SELECT COUNT(*) FROM characters").Scan(&count)
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
			db.Exec("INSERT INTO characters (name, rarity, is_pickup) VALUES ($1, $2, $3)",
				c.name, c.rarity, c.isPickup)
		}
	}

	// ガチャテーブル
	bannersTable := `
	CREATE TABLE IF NOT EXISTS gacha_banners (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		cost INTEGER NOT NULL DEFAULT 300,
		prob_star5 INTEGER NOT NULL DEFAULT 6,
		prob_star4 INTEGER NOT NULL DEFAULT 51,
		star5_limit INTEGER NOT NULL DEFAULT 90,
		star4_limit INTEGER NOT NULL DEFAULT 10,
		star5_pickup_prob INTEGER NOT NULL DEFAULT 100,
		pity_soft_start INTEGER NOT NULL DEFAULT 74,
		soft_pity_increment INTEGER NOT NULL DEFAULT 6
	);`
	if _, err := db.Exec(bannersTable); err != nil {
		return err
	}

	// ガチャがない場合デフォルトのガチャ設定を1つだけ自動作成する
	var bannerCount int
	db.QueryRow("SELECT COUNT(*) FROM gacha_banners").Scan(&bannerCount)
	if bannerCount == 0 {
		db.Exec(`INSERT INTO gacha_banners (title) VALUES ('恒常ガチャ')`)
	}

	// ピックアップテーブル（中間テーブル）
	pickupsTable := `
	CREATE TABLE IF NOT EXISTS banner_pickups (
		banner_id INTEGER REFERENCES gacha_banners(id) ON DELETE CASCADE,
		character_id INTEGER REFERENCES characters(id) ON DELETE CASCADE,
		PRIMARY KEY (banner_id, character_id)
	);`
	if _, err := db.Exec(pickupsTable); err != nil {
		return err
	}

	// 恒常テーブル（中間テーブル）
	constantTable := `
	CREATE TABLE IF NOT EXISTS constant_characters (
		PRIMARY character_id INTEGER REFERENCES characters(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(constantTable); err != nil {
		return err
	}

	// ガチャの履歴を保存するテーブルを作成
	historyTable := `
	CREATE TABLE IF NOT EXISTS history (
		id SERIAL PRIMARY KEY,
		uid TEXT REFERENCES users(uid) ON DELETE CASCADE,
		rarity TEXT NOT NULL,
		character TEXT NOT NULL
	);`
	if _, err = db.Exec(historyTable); err != nil {
		return err
	}

	// 決済の注文状態を管理するテーブルを作成
	ordersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		order_id TEXT PRIMARY KEY,
		uid TEXT,
		amount INTEGER,
		status TEXT   /* 'pending'(未払い) または 'paid'(支払い済み) */
	);`
	_, err = db.Exec(ordersTable)

	return err
}
