package main // エントリーポイント

// ライブラリのインポート
import (
	"database/sql"
	"encoding/json" // JSONのエンコード/デコードに使用
	core "gacha-core"
	"math/rand/v2" // 乱数 ガチャ用 (高速)
	"net/http"
)

// 定数の定義
const (
	// cookieの日数
	cookieDays = 30           // セッションIDを保存するCookieの有効期限（日数）
	oneDay     = 24 * 60 * 60 // 1日の秒数（CookieのMaxAgeに使用）
)

// ガチャの処理を行う関数
func (app *PublicApp) gachaHandler(w http.ResponseWriter, r *http.Request) {
	// CookieからユーザーIDを取得
	uid, err := getSession(r)
	if err != nil {
		http.Error(w, "ログインしてください", http.StatusUnauthorized)
		return
	}

	// ユーザーIDからユーザーデータを取得
	user := getUserData(app.db, uid)

	//ガチャバナーの取得
	gachaBanner := getGachaBanner(app.db, "恒常ガチャ")

	// 石の所持数をチェックして、足りない場合はエラーを返す
	if user.Stones < gachaBanner.Cost {
		http.Error(w, "石が足りません！", http.StatusBadRequest)
		return
	}

	// ガチャの結果を判定する関数を呼び出して、結果を取得
	result := gachaJudgment(app.db, user, "恒常ガチャ")

	// DB保存
	err = saveGachaResultTx(app.db, uid, user, []core.GachaResult{result}, gachaBanner.Cost)
	if err != nil {
		http.Error(w, "サーバーエラーが発生しました", http.StatusInternalServerError)
		return
	}

	// 石を消費
	user.Stones -= gachaBanner.Cost

	// 履歴に追加 (50件を超えていたら、一番古い要素を切り捨てる)
	user.GachaHistory = append(user.GachaHistory, result)
	if len(user.GachaHistory) > 50 {
		// インデックス1から最後までを残す
		user.GachaHistory = user.GachaHistory[1:]
	}

	// レスポンス作成
	sendGachaResponse(w, app.db, []core.GachaResult{result}, user)
}

// 10連ガチャの処理を行う関数
func (app *PublicApp) gacha10Handler(w http.ResponseWriter, r *http.Request) {
	// CookieからユーザーIDを取得
	uid, err := getSession(r)
	if err != nil {
		http.Error(w, "ログインしてください", http.StatusUnauthorized)
		return
	}

	// ユーザーIDからユーザーデータを取得
	user := getUserData(app.db, uid)

	//ガチャバナーの取得
	gachaBanner := getGachaBanner(app.db, "恒常ガチャ")

	// 石の所持数をチェックして、足りない場合はエラーを返す
	if user.Stones < gachaBanner.Cost*10 {
		http.Error(w, "石が足りません！", http.StatusBadRequest)
		return
	}

	var results []core.GachaResult
	for i := 0; i < 10; i++ {
		// ガチャの結果を判定する関数を呼び出して、結果を取得して、resultsの配列に追加
		result := gachaJudgment(app.db, user, "恒常ガチャ")
		results = append(results, result)

		// 履歴に追加 (50件を超えていたら、一番古い要素を切り捨てる)
		user.GachaHistory = append(user.GachaHistory, result)
		if len(user.GachaHistory) > 50 {
			// インデックス1から最後までを残す
			user.GachaHistory = user.GachaHistory[1:]
		}
	}

	// DB保存
	err = saveGachaResultTx(app.db, uid, user, results, gachaBanner.Cost*10)
	if err != nil {
		http.Error(w, "サーバーエラーが発生しました", http.StatusInternalServerError)
		return
	}

	// 石を消費
	user.Stones -= gachaBanner.Cost * 10

	// レスポンス作成
	sendGachaResponse(w, app.db, results, user)
}

// 天井カウンターを返すハンドラー
func (app *PublicApp) limitHandler(w http.ResponseWriter, r *http.Request) {
	// CookieからユーザーIDを取得
	uid, err := getSession(r)
	if err != nil {
		http.Error(w, "ログインしてください", http.StatusUnauthorized)
		return
	}

	// ユーザーIDからユーザーデータを取得
	user := getUserData(app.db, uid)

	//ガチャバナーの取得
	gachaBanner := getGachaBanner(app.db, "恒常ガチャ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]int{
		"star4LimitCounter": gachaBanner.Star4Limit - user.Star4LimitCounter,
		"star5LimitCounter": gachaBanner.Star5Limit - user.Star5LimitCounter,
		"stones":            user.Stones,
	})
}

// 履歴を返すハンドラー
func (app *PublicApp) historyHandler(w http.ResponseWriter, r *http.Request) {
	// CookieからユーザーIDを取得
	uid, err := getSession(r)
	if err != nil {
		http.Error(w, "ログインしてください", http.StatusUnauthorized)
		return
	}

	// ユーザーIDからユーザーデータを取得
	user := getUserData(app.db, uid)

	// 履歴が空の場合は、空の配列を返す
	if len(user.GachaHistory) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode([]core.GachaResult{})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(user.GachaHistory)
}

// 共通のレスポンス送信処理
func sendGachaResponse(w http.ResponseWriter, db *sql.DB, results []core.GachaResult, user *core.UserData) {
	//ガチャバナーの取得
	gachaBanner := getGachaBanner(db, "恒常ガチャ")

	response := core.GachaResponse{
		Results:   results,
		Pity5Star: gachaBanner.Star5Limit - user.Star5LimitCounter, // あと何回か
		Pity4Star: gachaBanner.Star4Limit - user.Star4LimitCounter,
		Stones:    user.Stones, // 所持石数
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

// ガチャの結果を判定する関数
func gachaJudgment(db *sql.DB, user *core.UserData, bannerTitle string) core.GachaResult {
	// カウンターをインクリメント
	user.Star4LimitCounter++ // 星4以上が出るまでのカウンター
	user.Star5LimitCounter++ // 星5が出るまでのカウンター

	//ガチャバナーの取得
	gachaBanner := getGachaBanner(db, bannerTitle)

	// ピックアップキャラの取得
	pickupCharacters := getPickupCharacters(db, bannerTitle)

	// 0〜999の乱数を生成
	roll := rand.IntN(1000)

	// 星5の当たる確率（6/1000 = 0.6%）
	star5Prob := gachaBanner.ProbBaseStar5

	// ソフトピティの確率上昇の判定
	if user.Star5LimitCounter >= gachaBanner.PitySoftStart {
		// 74連目以降は、6%ずつ確率が上昇
		star5Prob += gachaBanner.SoftPityIncrement * (user.Star5LimitCounter - (gachaBanner.PitySoftStart - 1))
	}

	// 確率の判定
	if roll < star5Prob || user.Star5LimitCounter >= gachaBanner.Star5Limit {
		// 0.6%の確率で星5 （もしくは、天井カウンターが90連目の場合は強制的に星5）
		user.Star4LimitCounter = 0 // カウンターをリセット
		user.Star5LimitCounter = 0 // カウンターをリセット

		// ピックアップキャラクターの当選判定を行う関数を呼び出す
		return pickupJudgment(user, gachaBanner.Star5PickupProb, pickupCharacters.Star5, getConstantCharacters(db))
	} else if roll < (star5Prob+gachaBanner.ProbBaseStar4) || user.Star4LimitCounter >= gachaBanner.Star4Limit {
		// 5.1%の確率で星4 （もしくは、天井カウンターが10連目の場合は強制的に星4）
		user.Star4LimitCounter = 0 // カウンターをリセット

		randomIndex := rand.IntN(len(pickupCharacters.Star4)) // ピックアップ星4キャラクターの中からランダムに選ぶ
		return core.GachaResult{Character: pickupCharacters.Star4[randomIndex]}
	} else {
		// 94.3%の確率で星3
		star3 := getStar3Characters(db)      // 星3キャラクターのリストをDBから取得
		randomIndex := rand.IntN(len(star3)) // 星3キャラクターの中からランダムに選ぶ
		return core.GachaResult{Character: star3[randomIndex]}
	}
}

// ピックアップキャラクターの当選判定を行う関数
func pickupJudgment(user *core.UserData, pickupProb int, pickupCharacters []core.Character, constantCharacters []core.Character) core.GachaResult {
	// ピックアップキャラクターが確定している場合は、ピックアップキャラクターを返す
	if user.IsNextPickupGuaranteed {
		user.IsNextPickupGuaranteed = false // フラグをリセット

		randomIndex := rand.IntN(len(pickupCharacters))                   // ピックアップ星5キャラクターの中からランダムに選ぶ
		return core.GachaResult{Character: pickupCharacters[randomIndex]} // ピックアップキャラクターの中から1体を返す
	} else {
		// ピックアップキャラクターが確定していない場合は、50%の確率でピックアップキャラクター、50%の確率ですり抜けキャラクターを返す
		if rand.IntN(100) < pickupProb {
			randomIndex := rand.IntN(len(pickupCharacters)) // ピックアップ星5キャラクターの中からランダムに選ぶ
			return core.GachaResult{Character: pickupCharacters[randomIndex]}
		} else {
			user.IsNextPickupGuaranteed = true // 次のガチャでピックアップキャラクターが確定するようにフラグをセット

			randomIndex := rand.IntN(len(constantCharacters)) // すり抜けキャラクターの中からランダムに選ぶ
			return core.GachaResult{Character: constantCharacters[randomIndex]}
		}
	}
}
