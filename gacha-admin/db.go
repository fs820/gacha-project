package main // エントリーポイント

import (
	"database/sql"
	"fmt"
	core "gacha-core"
	"log"
)

// DBから現在のバナーの配列を取得する関数
func getBanners(db *sql.DB) []core.GachaBanner {
	banners := []core.GachaBanner{}

	// DBから検索
	rows, err := db.Query(`
	    SELECT id title cost 
	    prob_star5, prob_star4, star5_limit, star4_limit, star5_pickup_prob, 
		pity_soft_start, soft_pity_increment
		FROM gacha_banners`)
	if err != nil {
		log.Println("キャラクター取得エラー:", err)
		return banners
	}
	defer rows.Close()

	// 取得したデータを構造体に格納
	for rows.Next() {
		var banner core.GachaBanner
		rows.Scan(&banner.ID, &banner.Title, &banner.Cost,
			&banner.ProbBaseStar5, &banner.ProbBaseStar4,
			&banner.Star5Limit, &banner.Star4Limit, &banner.Star5PickupProb,
			&banner.PitySoftStart, &banner.SoftPityIncrement)
		banners = append(banners, banner)
	}

	return banners
}

// 新しいバナーをDBに挿入する関数
func insertBanner(db *sql.DB, banner core.GachaBanner) error {
	// バナーを追加する
	_, err := db.Exec(`
	    INSERT INTO gacha_banners (id title cost 
	    prob_star5 prob_star4 star5_limit star4_limit star5_pickup_prob 
		pity_soft_start soft_pity_increment)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		`, banner.ID, banner.Title, banner.Cost,
		banner.ProbBaseStar5, banner.ProbBaseStar4,
		banner.Star5Limit, banner.Star4Limit, banner.Star5PickupProb,
		banner.PitySoftStart, banner.SoftPityIncrement)
	return err
}

// 指定したバナーの情報を書き換える関数
func changeBannersInfo(db *sql.DB, banner core.GachaBanner) error {
	// DBから検索
	_, err := db.Exec(`
	    UPDATE gacha_banners SET
	    title = $2, cost = $3, prob_star5 = $4, prob_star4 = $5,
		star5_limit = $6, star4_limit =$7, star5_pickup_prob = $8,
		pity_soft_start = $9, soft_pity_increment = $10
		WHERE id = $1
		`, banner.ID, banner.Title, banner.Cost,
		banner.ProbBaseStar5, banner.ProbBaseStar4,
		banner.Star5Limit, banner.Star4Limit, banner.Star5PickupProb,
		banner.PitySoftStart, banner.SoftPityIncrement)
	return err
}

// DBから現在のキャラクターの配列を取得する関数
func getCharacters(db *sql.DB) []core.Character {
	chars := []core.Character{}

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
	// キャラクターを追加する
	_, err := db.Exec("INSERT INTO characters (name, rarity) VALUES ($1, $2)",
		character.Name, character.Rarity)
	return err
}

// 恒常キャラを設定する関数
func changeConstantCharacter(db *sql.DB, constantCharacterIDs []int) error {
	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 恒常テーブルを空にする
	_, err = tx.Exec("DELETE FROM constant_characters")
	if err != nil {
		tx.Rollback()
		return err
	}

	// 新しいIDを追加する
	for _, ID := range constantCharacterIDs {
		_, err = tx.Exec(`
			INSERT INTO constant_characters character_id VALUES ($1)
		`, ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ピックアップの登録に失敗しました")
		}
	}

	return tx.Commit()
}

// DBから現在のキャラクターの配列を取得する関数
func getConstantCharactersID(db *sql.DB) []int {
	ids := []int{}

	// DBから検索
	rows, err := db.Query("SELECT character_id FROM constant_characters")
	if err != nil {
		log.Println("恒常取得エラー:", err)
		return ids
	}
	defer rows.Close()

	// 取得したデータを構造体に格納
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids
}

// DBから現在のキャラクターの配列を取得する関数
func getPickupCharactersID(db *sql.DB, bannerID int) []int {
	ids := []int{}

	// DBから検索
	rows, err := db.Query(`
	    SELECT character_id FROM banner_pickups WHERE banner_id = $1
		`, bannerID)
	if err != nil {
		log.Println("ピックアップ取得エラー:", err)
		return ids
	}
	defer rows.Close()

	// 取得したデータを構造体に格納
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids
}

// 指定したキャラクターをピックアップに設定する関数
func changePickupCharacter(db *sql.DB, bannerID int, pickupCharacters core.PickupCharacters) error {
	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 指定のガチャからピックアップをすべて解除する (gacha_bannersテーブルのtitleからガチャのidを探して指定)
	_, err = tx.Exec(`
		DELETE FROM banner_pickups WHERE banner_id = $1
	`, bannerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 選ばれたキャラクターたちをピックアップにする
	for _, ID := range pickupCharacters.Star5ID {
		// charactersテーブルから名前でIDを検索し、banner_idと一緒に登録する
		_, err = tx.Exec(`
			INSERT INTO banner_pickups (banner_id, character_id) VALUES ($1, $2)
		`, bannerID, ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ピックアップの登録に失敗しました")
		}
	}
	for _, ID := range pickupCharacters.Star4ID {
		// charactersテーブルから名前でIDを検索し、banner_idと一緒に登録する
		_, err = tx.Exec(`
			INSERT INTO banner_pickups (banner_id, character_id) VALUES ($1, $2)
		`, bannerID, ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ピックアップの登録に失敗しました")
		}
	}

	return tx.Commit()
}

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
