package main // エントリーポイント

// ライブラリのインポート
import (
	"encoding/json"
	"fmt"
	core "gacha-core"
	"net/http" // HTTPサーバーの構築に使用
)

// キャラクター追加のリクエスト型
type InsertCharacterRequest struct {
	Rarity string `json:"rarity"`
	Name   string `json:"name"`
}

// 石の付与のリクエスト型
type AddStoneRequest struct {
	UID    string `json:"uid"`
	Amount int    `json:"amount"`
}

// ピックアップ変更のリクエスト型
type UpdatePickupRequest struct {
	BannerTitle string   `json:"banner_title"`
	Star5Names  []string `json:"star5_names"`
	Star4Names  []string `json:"star4_names"`
}

// 管理者専用：すべての履歴を削除するエンドポイント
func (app *AdminApp) adminDeleteHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// 履歴テーブルのデータをすべて削除
	err := cleanupHistory(app.db)
	if err != nil {
		http.Error(w, "データベースの削除に失敗しました", http.StatusInternalServerError)
		return
	}

	// 成功メッセージ
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("すべてのガチャ履歴を正常に削除しました！"))
}

// 管理者専用：指定したユーザーに石を付与するエンドポイント
func (app *AdminApp) adminAddStonesHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// リクエストの読み込み
	var req AddStoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "不正なデータ形式です", http.StatusBadRequest)
		return
	}

	// 石を追加
	err := addStones(app.db, req.UID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 成功メッセージ (fmt.Sprintf を使って文字列の中に変数を埋め込む)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("ユーザー[%s]に石を%d個追加しました！", req.UID, req.Amount)))
}

// 管理者専用：キャラクター情報を追加するエンドポイント
func (app *AdminApp) adminInsertCharacterHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// リクエストの読み込み
	var req InsertCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "不正なデータ形式です", http.StatusBadRequest)
		return
	}

	// キャラクターを作る
	char := core.Character{
		Rarity: req.Rarity,
		Name:   req.Name,
	}

	// データベースの関数を呼び出して、指定したキャラクターを挿入
	err := insertCharacter(app.db, char)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("キャラクターが正常に追加されました！"))
}

// 管理者専用：ピックアップキャラクターを変更するエンドポイント
func (app *AdminApp) adminUpdatePickupHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// リクエストの読み込み
	var req UpdatePickupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "不正なデータ形式です", http.StatusBadRequest)
		return
	}

	// ピックアップキャラクターリストを作る
	var pickupCharacters core.PickupCharacters
	for _, name := range req.Star5Names {
		char := core.Character{
			Rarity: "星5",
			Name:   name,
		}
		pickupCharacters.Star5 = append(pickupCharacters.Star5, char)
	}
	for _, name := range req.Star4Names {
		char := core.Character{
			Rarity: "星4",
			Name:   name,
		}
		pickupCharacters.Star4 = append(pickupCharacters.Star4, char)
	}

	// データベースの関数を呼び出して、指定したキャラクターをピックアップに設定
	err := changePickupCharacter(app.db, req.BannerTitle, pickupCharacters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// メッセージ
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, name := range req.Star5Names {
		w.Write([]byte(fmt.Sprintf("星5ピックアップキャラクターを [%s] に更新しました！\n", name)))
	}
	for _, name := range req.Star4Names {
		w.Write([]byte(fmt.Sprintf("星4ピックアップキャラクターを [%s] に更新しました！\n", name)))
	}
}

// 管理者専用：キャラクター情報を取得するエンドポイント
func (app *AdminApp) adminGetCharacterHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// データベースの関数を呼び出して、指定したキャラクターの情報を取得
	characters := getCharacters(app.db)

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
