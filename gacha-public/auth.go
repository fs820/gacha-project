package main

// ライブラリのインポート
import (
	"crypto/rand" // 乱数 暗号用 (安全)
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// 登録情報構造体
type AuthRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RecaptchaToken string `json:"recaptcha_token"`
}

// セッションIDを生成する関数
func generateSessionID() string {
	b := make([]byte, 16)        // 16バイトのランダムなデータを格納するためのバイトスライスを作成
	rand.Read(b)                 // バイトスライスにランダムなデータを埋める
	return hex.EncodeToString(b) // バイトスライスを16進数の文字列に変換して返す
}

// CookieからセッションIDを取得する関数
func getSession(r *http.Request) (string, error) {
	// リクエストから "session_id" という名前のCookieを探す
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Googleのサーバーに通信してトークンを検証する関数
func verifyRecaptcha(token string) bool {
	// .env ファイルから reCAPTCHA のシークレットキーを取得
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		// 検証失敗
		return false
	}

	// GoogleのAPIのURL
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	// POSTするデータ（x-www-form-urlencoded 形式）
	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", token)

	// GoogleのAPIにPOSTリクエストを送る
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return false // 通信エラー
	}
	defer resp.Body.Close()

	// Googleからの返事（JSON）を読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// 返事のJSONをパースする
	var result struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	// Googleが人間ですと判定したか
	return result.Success
}

// ログイン状態を確認するハンドラー
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	// CookieからユーザーIDを取得
	_, err := getSession(r)
	if err != nil {
		// Cookieが無い、または期限切れの場合はログインが必要
		http.Error(w, "未ログイン", http.StatusUnauthorized)
		return
	}

	// 取得できればOK (ログイン状態)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// 新規登録ハンドラー
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "POSTリクエストのみ許可されています", http.StatusMethodNotAllowed)
		return
	}

	// 受け取った登録情報をデコード
	var req AuthRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Bot対策：reCAPTCHAトークンの検証
	if !verifyRecaptcha(req.RecaptchaToken) {
		http.Error(w, "Botによるアクセスの可能性があります", http.StatusForbidden)
		return
	}

	// 空チェック
	if req.Username == "" || req.Password == "" {
		http.Error(w, "ユーザー名とパスワードを入力してください", http.StatusBadRequest)
		return
	}

	// パスワードを暗号化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "パスワードの暗号化に失敗しました", http.StatusInternalServerError)
		return
	}

	// 新しいユーザーIDを生成
	uid := generateSessionID()

	// DB登録
	err = insertUser(uid, req.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "そのユーザー名は既に使われています", http.StatusConflict)
		return
	}

	// 成功メッセージ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ユーザー登録が完了しました！ログインしてください。"))
}

// ログインハンドラー
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "POSTリクエストのみ許可されています", http.StatusMethodNotAllowed)
		return
	}

	// 受け取ったログイン情報をデコード
	var req AuthRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Bot対策：reCAPTCHAトークンの検証
	if !verifyRecaptcha(req.RecaptchaToken) {
		http.Error(w, "Botによるアクセスの可能性があります", http.StatusForbidden)
		return
	}

	uid, hashedPassword, err := findUser(req.Username)
	if err != nil {
		http.Error(w, "ユーザー名またはパスワードが間違っています", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "ユーザー名またはパスワードが間違っています", http.StatusUnauthorized)
		return
	}

	// Cookieをログインユーザーのuidで上書きする
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    uid,
		Path:     "/",                 // サイト内の全ページでこのCookieを有効にする
		HttpOnly: true,                // JavaScriptからCookieを盗まれるのを防ぐセキュリティ設定
		MaxAge:   oneDay * cookieDays, // 有効期限（秒数）
	}
	http.SetCookie(w, cookie)

	w.Write([]byte("ログインに成功しました！"))
}
