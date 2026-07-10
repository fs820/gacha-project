package main // エントリーポイント

// ライブラリのインポート
import (
	"encoding/json"
	"fmt"
	"net/http" // HTTPサーバーの構築に使用
	"strconv"
	"strings"
)

// キャラクター情報を表す構造体
type Character struct {
	Name     string `json:"name"`
	Rarity   string `json:"rarity"`
	IsPickup bool   `json:"isPickup"`
}

const PASSWORD = "Bearer supersecret"

// 管理者専用：すべての履歴を削除するエンドポイント
func adminDeleteHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// パスワードチェック
	authHeader := r.Header.Get("Authorization")
	if authHeader != PASSWORD {
		http.Error(w, "権限がありません (Unauthorized)", http.StatusUnauthorized)
		return
	}

	// 履歴テーブルのデータをすべて削除
	err := cleanupHistory()
	if err != nil {
		http.Error(w, "データベースの削除に失敗しました", http.StatusInternalServerError)
		return
	}

	// 成功メッセージ
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("すべてのガチャ履歴を正常に削除しました！"))
}

// 管理者専用：指定したユーザーに石を付与するエンドポイント
func adminAddStonesHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// パスワードチェック
	authHeader := r.Header.Get("Authorization")
	if authHeader != PASSWORD {
		http.Error(w, "権限がありません (Unauthorized)", http.StatusUnauthorized)
		return
	}

	// クエリパラメータの取得
	targetUID := r.URL.Query().Get("uid")
	amountStr := r.URL.Query().Get("amount")
	if targetUID == "" || amountStr == "" {
		http.Error(w, "uidとamountを指定してください。 例: ?pass=supersecret&uid=xxx&amount=1000", http.StatusBadRequest)
		return
	}

	// 文字列のamountを整数(int)に変換
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "amountは数字で指定してください", http.StatusBadRequest)
		return
	}

	// 石を追加
	err = addStones(targetUID, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 成功メッセージ (fmt.Sprintf を使って文字列の中に変数を埋め込む)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("ユーザー[%s]に石を%d個追加しました！", targetUID, amount)))
}

// 管理者専用：キャラクター情報を追加するエンドポイント
func adminInsertCharacterHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// パスワードチェック
	authHeader := r.Header.Get("Authorization")
	if authHeader != PASSWORD {
		http.Error(w, "権限がありません (Unauthorized)", http.StatusUnauthorized)
		return
	}

	// URLから変更したいキャラクターの名前を取得（例: ?name=アテナ）
	targetNamesStr := r.URL.Query().Get("name")
	if targetNamesStr == "" {
		http.Error(w, "nameを指定してください。 例: ?rarity=星5&name=アテナ", http.StatusBadRequest)
		return
	}

	// URLから変更したいレアリティを取得（例: ?rarity=星5）
	rarity := r.URL.Query().Get("rarity")
	if rarity == "" {
		http.Error(w, "rarityを指定してください。 例: ?rarity=星5&name=アテナ", http.StatusBadRequest)
		return
	}

	// 新しいキャラクター情報を作成
	newCharacter := Character{
		Name:   targetNamesStr,
		Rarity: rarity,
	}

	// データベースの関数を呼び出して、指定したキャラクターを挿入
	err := insertCharacter(newCharacter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("キャラクターが正常に追加されました！"))
}

// 管理者専用：ピックアップキャラクターを変更するエンドポイント
func adminUpdatePickupHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// パスワードチェック
	authHeader := r.Header.Get("Authorization")
	if authHeader != PASSWORD {
		http.Error(w, "権限がありません (Unauthorized)", http.StatusUnauthorized)
		return
	}

	// URLから変更したいレアリティを取得（例: ?rarity=星5）
	rarity := r.URL.Query().Get("rarity")
	if rarity == "" {
		http.Error(w, "rarityを指定してください。 例: ?rarity=星5&name=アテナ", http.StatusBadRequest)
		return
	}

	// URLから変更したいキャラクターの名前を取得（例: ?name=アテナ）
	targetNamesStr := r.URL.Query().Get("name")
	if targetNamesStr == "" {
		http.Error(w, "nameを指定してください。 例: ?rarity=星5&name=アテナ", http.StatusBadRequest)
		return
	}
	targetNames := strings.Split(targetNamesStr, ",")
	for i := range targetNames {
		if targetNames[i] == "" {
			http.Error(w, "nameを指定してください。 例: ?rarity=星5&name=アテナ", http.StatusBadRequest)
			return
		}
	}

	// データベースの関数を呼び出して、指定したキャラクターをピックアップに設定
	err := changePickupCharacter(rarity, targetNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, name := range targetNames {
		w.Write([]byte(fmt.Sprintf("%sピックアップキャラクターを [%s] に更新しました！\n", rarity, name)))
	}
}

// 管理者専用：キャラクター情報を取得するエンドポイント
func adminGetCharacterHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないリクエスト方法です (Method Not Allowed)", http.StatusMethodNotAllowed)
		return
	}

	// パスワードチェック
	authHeader := r.Header.Get("Authorization")
	if authHeader != PASSWORD {
		http.Error(w, "権限がありません (Unauthorized)", http.StatusUnauthorized)
		return
	}

	// データベースの関数を呼び出して、指定したキャラクターの情報を取得
	characters := getCharacters()

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
