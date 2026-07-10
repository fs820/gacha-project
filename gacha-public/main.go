package main // エントリーポイント

// ライブラリのインポート
import (
	"fmt" // フォーマット用 (文字列の整形など)
	"log"
	"net/http" // HTTPサーバーの構築に使用

	"github.com/joho/godotenv" // .env ファイルを読み込むためのライブラリ
)

// メイン関数
func main() {
	// 【新規追加】サーバー起動時に .env ファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: .env ファイルが見つかりません。環境変数が直接設定されているか確認してください。")
	}

	// データベースの初期化
	initDB()

	// "static"フォルダの中身（HTML, CSS, JS）を、そのままブラウザに公開する設定
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// ガチャのエンドポイントを設定
	http.HandleFunc("/gacha", gachaHandler)     // 単発ガチャのエンドポイント /gacha
	http.HandleFunc("/gacha10", gacha10Handler) // 10連ガチャのエンドポイント /gacha10

	// 履歴だけを取得するエンドポイント
	http.HandleFunc("/history", historyHandler)
	// 天井カウンターを取得するエンドポイント
	http.HandleFunc("/limit", limitHandler)

	// 認証用エンドポイント
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/check_auth", checkAuthHandler)

	// 決済用エンドポイント
	http.HandleFunc("/checkout", checkoutHandler)
	http.HandleFunc("/webhook/payment", webhookHandler)

	// サーバー起動のメッセージを表示
	fmt.Println("サーバーを起動しました！ ブラウザで http://localhost:8080 にアクセスしてください。")
	fmt.Println("終了するにはターミナルで Ctrl + C を押します。")

	// ポート8080でサーバーを起動（ゲームのメインループのように、ここでアクセスを待ち続けます）
	http.ListenAndServe(":8080", nil)
}
