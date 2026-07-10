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
	// サーバー起動時に環境変数(.envやRenderの環境変数)を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: .env ファイルが見つかりません。環境変数が直接設定されているか確認してください。")
	}

	// データベースの初期化
	initDB()

	// "static"フォルダの中身（HTML, CSS, JS）を、そのままブラウザに公開する設定
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// 管理者用エンドポイント
	http.HandleFunc("/admin/delete_history", adminDeleteHistoryHandler)
	http.HandleFunc("/admin/add_stones", adminAddStonesHandler)
	http.HandleFunc("/admin/insert_character", adminInsertCharacterHandler)
	http.HandleFunc("/admin/update_pickup", adminUpdatePickupHandler)
	http.HandleFunc("/admin/get_character", adminGetCharacterHandler)

	// サーバー起動のメッセージを表示
	fmt.Println("サーバーを起動しました！ ブラウザで http://localhost:8081 にアクセスしてください。")
	fmt.Println("終了するにはターミナルで Ctrl + C を押します。")

	// ポート8081でサーバーを起動（ゲームのメインループのように、ここでアクセスを待ち続けます）
	http.ListenAndServe(":8081", nil)
}
