package main // エントリーポイント

// ライブラリのインポート
import (
	"encoding/json"
	"fmt"
	core "gacha-core"
	"net/http" // HTTPサーバーの構築に使用
)

// 基本レスポンス
type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// バナーの追加のリクエスト型
type InsertBannerRequest struct {
	Title             string `json:"title"`
	Cost              int    `json:"cost"`
	ProbBaseStar5     int    `json:"probBaseStar5"`
	ProbBaseStar4     int    `json:"probBaseStar4"`
	Star5Limit        int    `json:"star5Limit"`
	Star4Limit        int    `json:"star4Limit"`
	Star5PickupProb   int    `json:"star5PickupProb"`
	PitySoftStart     int    `json:"pitySoftStart"`
	SoftPityIncrement int    `json:"softPityIncrement"`
}

// キャラクター追加のリクエスト型
type InsertCharacterRequest struct {
	Rarity string `json:"rarity"`
	Name   string `json:"name"`
}

// 恒常変更のリクエスト型
type UpdateConstantRequest struct {
	CharID []int `json:"char_id"`
}

// ピックアップ変更のリクエスト型
type UpdatePickupRequest struct {
	BannerID int   `json:"banner_id"`
	Star5ID  []int `json:"star5_id"`
	Star4ID  []int `json:"star4_id"`
}

// 石の付与のリクエスト型
type AddStoneRequest struct {
	UID    string `json:"uid"`
	Amount int    `json:"amount"`
}

// 成功メッセージ
func sendSuccessResponse(w http.ResponseWriter, message string) {
	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	apiResponse := ApiResponse{
		Success: true,
		Message: message,
	}
	json.NewEncoder(w).Encode(apiResponse)
}

// 管理者専用：バナー情報を取得するエンドポイント
func (app *AdminApp) adminGetBannerHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// データベースの関数を呼び出して、指定したキャラクターの情報を取得
	banners := getBanners(app.db)

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// 管理者専用：バナー情報を追加するエンドポイント
func (app *AdminApp) adminInsertBannerHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// リクエストの読み込み
	var req InsertBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "不正なデータ形式です", http.StatusBadRequest)
		return
	}

	// バナーを作る
	banner := core.GachaBanner{
		Title:             req.Title,
		Cost:              req.Cost,
		ProbBaseStar5:     req.ProbBaseStar5,
		ProbBaseStar4:     req.ProbBaseStar4,
		Star5Limit:        req.Star5Limit,
		Star4Limit:        req.Star4Limit,
		Star5PickupProb:   req.Star5PickupProb,
		PitySoftStart:     req.PitySoftStart,
		SoftPityIncrement: req.SoftPityIncrement,
	}

	// データベースの関数を呼び出して、指定したキャラクターを挿入
	err := insertBanner(app.db, banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON形式でレスポンスを返す
	sendSuccessResponse(w, "バナーが正常に追加されました！")
}

// 管理者専用：バナー情報を変更するエンドポイント
func (app *AdminApp) adminChangeBannerHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// リクエストの読み込み
	var req core.GachaBanner
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "不正なデータ形式です", http.StatusBadRequest)
		return
	}

	// データベースの関数を呼び出して、指定したキャラクターを挿入
	err := changeBannersInfo(app.db, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON形式でレスポンスを返す
	sendSuccessResponse(w, "バナーが正常に変更されました！")
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

	// JSON形式でレスポンスを返す
	sendSuccessResponse(w, "キャラクターが正常に追加されました！")
}

// 管理者専用：恒常キャラクターを変更するエンドポイント
func (app *AdminApp) adminUpdateConstantHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// リクエストの読み込み
	var req UpdateConstantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "不正なデータ形式です", http.StatusBadRequest)
		return
	}

	// ピックアップキャラクターリストを作る
	var constantCharacters []int
	for _, ID := range req.CharID {
		constantCharacters = append(constantCharacters, ID)
	}

	// データベースの関数を呼び出して、指定したキャラクターをピックアップに設定
	err := changeConstantCharacter(app.db, constantCharacters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON形式でレスポンスを返す
	for _, ID := range req.CharID {
		sendSuccessResponse(w, fmt.Sprintf("恒常キャラクターを [%d] に更新しました！\n", ID))
	}
}

// 管理者専用：ピックアップID情報を取得するエンドポイント
func (app *AdminApp) adminGetConstantIDHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// データベースの関数を呼び出して、指定したキャラクターの情報を取得
	ids := getConstantCharactersID(app.db)

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)
}

// 管理者専用：ピックアップID情報を取得するエンドポイント
func (app *AdminApp) adminGetPickupIDHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// リクエストの読み込み
	var bannerID int
	if err := json.NewDecoder(r.Body).Decode(&bannerID); err != nil {
		http.Error(w, "不正なデータ形式です", http.StatusBadRequest)
		return
	}

	// データベースの関数を呼び出して、指定したキャラクターの情報を取得
	ids := getPickupCharactersID(app.db, bannerID)

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)
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
	for _, ID := range req.Star5ID {
		pickupCharacters.Star5ID = append(pickupCharacters.Star5ID, ID)
	}
	for _, ID := range req.Star4ID {
		pickupCharacters.Star4ID = append(pickupCharacters.Star4ID, ID)
	}

	// データベースの関数を呼び出して、指定したキャラクターをピックアップに設定
	err := changePickupCharacter(app.db, req.BannerID, pickupCharacters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON形式でレスポンスを返す
	for _, ID := range req.Star5ID {
		sendSuccessResponse(w, fmt.Sprintf("星5ピックアップキャラクターを [%d] に更新しました！\n", ID))
	}
	for _, ID := range req.Star4ID {
		sendSuccessResponse(w, fmt.Sprintf("星4ピックアップキャラクターを [%d] に更新しました！\n", ID))
	}
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

	// JSON形式でレスポンスを返す
	sendSuccessResponse(w, fmt.Sprintf("ユーザー[%s]に石を%d個追加しました！", req.UID, req.Amount))
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

	// JSON形式でレスポンスを返す
	sendSuccessResponse(w, "すべてのガチャ履歴を正常に削除しました！")
}
