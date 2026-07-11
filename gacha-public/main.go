package main // エントリーポイント

// ライブラリのインポート
import (
	"database/sql"
	"fmt"             // フォーマット用 (文字列の整形など)
	core "gacha-core" // DBコア
	"log"             // ロギング
	"net/http"        // HTTPサーバーの構築

	"github.com/joho/godotenv" // .env ファイルを読み込むためのライブラリ
)

// 管理者用アプリケーション構造体
type PublicApp struct {
	db *sql.DB
}

// メイン関数
func main() {
	// 【新規追加】サーバー起動時に .env ファイルを読み込む
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

	// アプリケーションの構築
	app := &PublicApp{db: database}

	// "static"フォルダの中身（HTML, CSS, JS）を、そのままブラウザに公開する設定
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// ガチャのエンドポイントを設定
	http.HandleFunc("/gacha", app.gachaHandler)     // 単発ガチャのエンドポイント /gacha
	http.HandleFunc("/gacha10", app.gacha10Handler) // 10連ガチャのエンドポイント /gacha10

	// 履歴だけを取得するエンドポイント
	http.HandleFunc("/history", app.historyHandler)
	// 天井カウンターを取得するエンドポイント
	http.HandleFunc("/limit", app.limitHandler)

	// 認証用エンドポイント
	http.HandleFunc("/register", app.registerHandler)
	http.HandleFunc("/login", app.loginHandler)
	http.HandleFunc("/check_auth", app.checkAuthHandler)

	// 決済用エンドポイント
	http.HandleFunc("/checkout", app.checkoutHandler)
	http.HandleFunc("/webhook/payment", app.webhookHandler)

	// サーバー起動のメッセージを表示
	fmt.Println("サーバーを起動しました！ ブラウザで http://localhost:8080 にアクセスしてください。")
	fmt.Println("終了するにはターミナルで Ctrl + C を押します。")

	// ポート8080でサーバーを起動（ゲームのメインループのように、ここでアクセスを待ち続けます）
	http.ListenAndServe(":8080", nil)
}
