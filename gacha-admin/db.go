package main // エントリーポイント

import (
	"database/sql"
	"fmt"
	core "gacha-core"
	"log"
)

// ガチャ石を追加する関数
func addStones(db *sql.DB, uid string, stonesToAdd int) error {
	// 石を追加する
	_, err := db.Exec("UPDATE users SET stones = stones + $1 WHERE uid = $2", stonesToAdd, uid)
	return err
}

// 履歴テーブルのデータをすべて削除する関数
func cleanupHistory(db *sql.DB) error {
	// PostgreSQLの一括削除とIDリセット
	_, err := db.Exec("TRUNCATE TABLE history RESTART IDENTITY")
	return err
}

// DBから現在のキャラクターの配列を取得する関数
func getCharacters(db *sql.DB) []core.Character {
	var chars []core.Character

	// DBから検索
	rows, err := db.Query("SELECT name, rarity, is_pickup FROM characters")
	if err != nil {
		log.Println("キャラクター取得エラー:", err)
		return chars
	}
	defer rows.Close()

	// 取得したデータを構造体に格納
	for rows.Next() {
		var char core.Character
		rows.Scan(&char.Name, &char.Rarity, &char.IsPickup)
		chars = append(chars, char)
	}

	return chars
}

// 新しいキャラクターをDBに挿入する関数
func insertCharacter(db *sql.DB, Character core.Character) error {
	_, err := db.Exec("INSERT INTO characters (name, rarity, is_pickup) VALUES ($1, $2, $3)",
		Character.Name, Character.Rarity, Character.IsPickup)
	return err
}

// 指定したキャラクターをピックアップに設定する関数
func changePickupCharacter(db *sql.DB, rarity string, targetNames []string) error {
	// トランザクション開始
	tx, err := db.Begin()
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
