package main // エントリーポイント

// ライブラリのインポート
import (
	"database/sql" // データベース操作に使用
	"fmt"
	"log" // ロギングに使用
	"os"

	_ "github.com/lib/pq" // PostgreSQLドライバ (データベース接続)
)

var userDB *sql.DB

// データベースの初期化関数
func initDB() {
	var err error
	// .env や Render の環境変数からURLを取得
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	// PostgreSQLに接続
	userDB, err = sql.Open("postgres", dbURL)
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
	_, err = userDB.Exec(usersTable)
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
	_, err = userDB.Exec(historyTable)
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
	_, err = userDB.Exec(charactersTable)
	if err != nil {
		log.Fatal("charactersテーブル作成エラー:", err)
	}

	// もしキャラクターテーブルが空の場合初期化
	var count int
	userDB.QueryRow("SELECT COUNT(*) FROM characters").Scan(&count)
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
			userDB.Exec("INSERT INTO characters (name, rarity, is_pickup) VALUES ($1, $2, $3)",
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
	_, err = userDB.Exec(ordersTable)
	if err != nil {
		log.Fatal("ordersテーブル作成エラー:", err)
	}
}

// ガチャ石を追加する関数
func addStones(uid string, stonesToAdd int) error {
	// 石を追加する
	_, err := userDB.Exec("UPDATE users SET stones = stones + $1 WHERE uid = $2", stonesToAdd, uid)
	return err
}

// 履歴テーブルのデータをすべて削除する関数
func cleanupHistory() error {
	// PostgreSQLの一括削除とIDリセット
	_, err := userDB.Exec("TRUNCATE TABLE history RESTART IDENTITY")
	return err
}

// DBから現在のキャラクターの配列を取得する関数
func getCharacters() []Character {
	var chars []Character

	// DBから検索
	rows, err := userDB.Query("SELECT name, rarity, is_pickup FROM characters")
	if err != nil {
		log.Println("キャラクター取得エラー:", err)
		return chars
	}
	defer rows.Close()

	// 取得したデータを構造体に格納
	for rows.Next() {
		var char Character
		rows.Scan(&char.Name, &char.Rarity, &char.IsPickup)
		chars = append(chars, char)
	}

	return chars
}

// 新しいキャラクターをDBに挿入する関数
func insertCharacter(Character Character) error {
	_, err := userDB.Exec("INSERT INTO characters (name, rarity, is_pickup) VALUES ($1, $2, $3)",
		Character.Name, Character.Rarity, Character.IsPickup)
	return err
}

// 指定したキャラクターをピックアップに設定する関数
func changePickupCharacter(rarity string, targetNames []string) error {
	// トランザクション開始
	tx, err := userDB.Begin()
	if err != nil {
		return err
	}

	// 指定のレアリティのキャラクターの is_pickup を一旦 false (非ピックアップ) にリセットする
	_, err = tx.Exec("UPDATE characters SET is_pickup = false WHERE rarity = $1", rarity)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 指定されたキャラクターの is_pickup だけを true にする
	for _, targetName := range targetNames {
		res, err := tx.Exec("UPDATE characters SET is_pickup = true WHERE name = $1", targetName)
		if err != nil {
			tx.Rollback()
			return err
		}

		// もし指定した名前のキャラが存在しなかった場合
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("指定されたキャラクター名[%s]は存在しません", targetName)
		}
	}

	return tx.Commit()
}
