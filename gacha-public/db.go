package main // エントリーポイント

// ライブラリのインポート
import (
	"database/sql"    // データベース操作に使用
	"fmt"             // フォーマット
	core "gacha-core" // DBコア
	"log"             // ロギングに使用
)

// ユーザーIDからデータを取得する関数
func getUserData(db *sql.DB, uid string) *core.UserData {
	user := &core.UserData{}
	var isGuaranteed bool // PostgreSQLのBOOLEANを安全に受け取るための一時変数

	// カウンター情報の取得
	row := db.QueryRow("SELECT stones, star4_limit_counter, star5_limit_counter, is_next_pickup_guaranteed FROM users WHERE uid = $1", uid)
	err := row.Scan(&user.Stones, &user.Star4LimitCounter, &user.Star5LimitCounter, &isGuaranteed)
	if err == sql.ErrNoRows {
		// データが無い（新規ユーザー）の場合は、初期値をDBに登録
		db.Exec("INSERT INTO users (uid) VALUES ($1)", uid)
	}

	// ピックアップ保証の状態をUserDataに反映
	user.IsNextPickupGuaranteed = isGuaranteed

	// 履歴の取得新しいものを50件取得して、古い順に並び替えるs
	rows, err := db.Query("SELECT rarity, character FROM (SELECT id, rarity, character FROM history WHERE uid = $1 ORDER BY id DESC LIMIT 50) AS sub ORDER BY id ASC", uid)
	if err != nil {
		log.Println("履歴取得エラー:", err)
		return user // エラーが起きたらここで中断
	}
	defer rows.Close() // 使い終わったら必ず閉じる

	// 取得した履歴をUserDataのGachaHistoryに追加
	for rows.Next() {
		var res core.GachaResult
		var char core.Character
		rows.Scan(&char.Rarity, &char.Name)
		res.Character = char
		user.GachaHistory = append(user.GachaHistory, res)
	}

	return user
}

// DBに新規ユーザーを登録する関数 ※パスワードはhash済文字列
func insertUser(db *sql.DB, uid string, username string, hashedPassword string) error {
	// データベースに新しいユーザーを保存
	_, err := db.Exec("INSERT INTO users (uid, username, password_hash) VALUES ($1, $2, $3)",
		uid, username, hashedPassword)
	return err
}

// ユーザー名からDBを検索して uidとhash済パスワードを返す関数
func findUser(db *sql.DB, username string) (string, string, error) {
	var uid, hash string
	err := db.QueryRow("SELECT uid, password_hash FROM users WHERE username = $1", username).Scan(&uid, &hash)
	return uid, hash, err
}

// ガチャの結果を保存する関数 （トランザクション）
func saveGachaResultTx(db *sql.DB, uid string, user *core.UserData, results []core.GachaResult, cost int) error {
	// トランザクションの開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 石を消費してカウンターを進める
	_, err = tx.Exec("UPDATE users SET stones = stones - $1, star4_limit_counter = $2, star5_limit_counter = $3, is_next_pickup_guaranteed = $4 WHERE uid = $5",
		cost, user.Star4LimitCounter, user.Star5LimitCounter, user.IsNextPickupGuaranteed, uid)
	if err != nil {
		tx.Rollback() // エラーが起きたらロールバック
		return err
	}

	// ガチャの結果を履歴テーブルに保存
	for _, res := range results {
		_, err = tx.Exec("INSERT INTO history (uid, rarity, character) VALUES ($1, $2, $3)", uid, res.Character.Rarity, res.Character.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// コミットして確定
	return tx.Commit()
}

// DBから指定したレアリティとピックアップ条件に合うキャラクターの配列を取得する関数
func getCharacters(db *sql.DB, rarity string, isPickup bool) []core.Character {
	var chars []core.Character

	// DBから検索
	rows, err := db.Query("SELECT name FROM characters WHERE rarity = $1 AND is_pickup = $2", rarity, isPickup)
	if err != nil {
		log.Println("キャラクター取得エラー:", err)
		return chars
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		rows.Scan(&name)
		char := core.Character{
			Name:     name,
			Rarity:   rarity,
			IsPickup: isPickup,
		}
		chars = append(chars, char)
	}

	return chars
}

// ユーザーの石を購入するリクエストをDBに登録する関数
func registerOrder(db *sql.DB, orderID string, uid string, amount int) error {
	_, err := db.Exec("INSERT INTO orders (order_id, uid, amount, status) VALUES ($1, $2, $3, 'pending')", orderID, uid, amount)
	return err
}

// 決済会社が決済出来た時に呼ばれる、石を増やしえ決済を完了する関数 (トランザクション)
func completeOrderTx(db *sql.DB, orderID string) error {
	// トランザクション開始 (注文の完了、石の付与)
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 注文を完了状態にに更新する (すでに 完了している('paid')の場合は何もしない(二重決済禁止))
	res, err := tx.Exec("UPDATE orders SET status = 'paid' WHERE order_id = $1 AND status = 'pending'", orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 実際に更新された行数をチェック
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("無効な注文番号、または既に支払い済みの注文です")
	}

	// 購入したユーザーのIDと、付与する石の量を取得
	var uid string
	var amount int
	err = tx.QueryRow("SELECT uid, amount FROM orders WHERE order_id = $1", orderID).Scan(&uid, &amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 石を付与
	_, err = tx.Exec("UPDATE users SET stones = stones + $1 WHERE uid = $2", amount, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	// トランザクション終了
	return tx.Commit()
}
