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
	rows, err := db.Query("SELECT id name, rarity FROM characters")
	if err != nil {
		log.Println("キャラクター取得エラー:", err)
		return chars
	}
	defer rows.Close()

	// 取得したデータを構造体に格納
	for rows.Next() {
		var char core.Character
		rows.Scan(&char.ID, &char.Name, &char.Rarity)
		chars = append(chars, char)
	}

	return chars
}

// 新しいキャラクターをDBに挿入する関数
func insertCharacter(db *sql.DB, character core.Character) error {
	// キャラクターを追加してIDを取得
	var newID int
	err := db.QueryRow("INSERT INTO characters (name, rarity) VALUES ($1, $2) RETURNING id",
		character.Name, character.Rarity).Scan(&newID)
	if err != nil {
		return err
	}

	return nil
}

// 指定したキャラクターをピックアップに設定する関数
func changePickupCharacter(db *sql.DB, bannerTitle string, pickupCharacters core.PickupCharacters) error {
	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 指定のガチャからピックアップをすべて解除する (gacha_bannersテーブルのtitleからガチャのidを探して指定)
	_, err = tx.Exec(`
		DELETE FROM banner_pickups 
		WHERE banner_id = (SELECT id FROM gacha_banners WHERE title = $1)
	`, bannerTitle)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 選ばれたキャラクターたちをピックアップにする
	for _, character := range pickupCharacters.Star5 {
		// charactersテーブルから名前でIDを検索し、banner_id=1と一緒に登録する
		_, err = tx.Exec(`
			INSERT INTO banner_pickups (banner_id, character_id)
			SELECT (SELECT id FROM gacha_banners WHERE title = $1), (id FROM characters WHERE name = $2)
		`, bannerTitle, character.Name)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ピックアップの登録に失敗しました")
		}
	}
	for _, character := range pickupCharacters.Star4 {
		// charactersテーブルから名前でIDを検索し、banner_id=1と一緒に登録する
		_, err = tx.Exec(`
			INSERT INTO banner_pickups (banner_id, character_id)
			SELECT (SELECT id FROM gacha_banners WHERE title = $1), (id FROM characters WHERE name = $2)
		`, bannerTitle, character.Name)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ピックアップの登録に失敗しました")
		}
	}

	return tx.Commit()
}
