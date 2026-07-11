package main // エントリーポイント

// ライブラリのインポート
import (
	"database/sql"
	"fmt" // フォーマット用 (文字列の整形など)
	core "gacha-core"
	"log"
	"net/http" // HTTPサーバーの構築に使用
	"os"

	"github.com/joho/godotenv" // .env ファイルを読み込むためのライブラリ
)

// 管理者用アプリケーション構造体
type AdminApp struct {
	db *sql.DB
}

// Basic認証用のミドルウェア関数を作る
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		password := os.Getenv("PASSWORD")
		// ユーザー名とパスワードを判定
		if !ok || user != "admin" || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// メイン関数
func main() {
	// サーバー起動時に環境変数(.envやRenderの環境変数)を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: .env ファイルが見つかりません。環境変数が直接設定されているか確認してください。")
	}

	// DBに接続する
	database, err := core.NewDBConnection()
	if err != nil {
		log.Fatal("DB接続エラー:", err)
	}
	defer database.Close()

	// データベースの初期化
	core.InitSchema(database)

	// アプリケーションの構築
	app := &AdminApp{db: database}

	// "static"フォルダの中身（HTML, CSS, JS）を、そのままブラウザに公開する設定
	fs := http.FileServer(http.Dir("static"))

	// 静的ファイルを提供する (Basic認証)
	http.Handle("/", basicAuth(fs))

	// 管理者用エンドポイント (Basic認証)
	http.Handle("/admin/delete_history", basicAuth(http.HandlerFunc(app.adminDeleteHistoryHandler)))     // 履歴の削除
	http.Handle("/admin/add_stones", basicAuth(http.HandlerFunc(app.adminAddStonesHandler)))             // 石の付与
	http.Handle("/admin/insert_character", basicAuth(http.HandlerFunc(app.adminInsertCharacterHandler))) // キャラクター追加
	http.Handle("/admin/update_pickup", basicAuth(http.HandlerFunc(app.adminUpdatePickupHandler)))       // ピックアップ変更
	http.Handle("/admin/get_character", basicAuth(http.HandlerFunc(app.adminGetCharacterHandler)))       // キャラクター取得

	// サーバー起動のメッセージを表示
	fmt.Println("サーバーを起動しました！ ブラウザで http://localhost:8081 にアクセスしてください。")
	fmt.Println("終了するにはターミナルで Ctrl + C を押します。")

	// ポート8081でサーバーを起動（ゲームのメインループのように、ここでアクセスを待ち続けます）
	http.ListenAndServe(":8081", nil)
}
