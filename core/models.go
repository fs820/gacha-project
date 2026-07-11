package core

// ユーザデータ
type UserData struct {
	Stones                 int
	Star4LimitCounter      int
	Star5LimitCounter      int
	IsNextPickupGuaranteed bool
	GachaHistory           []GachaResult
}

// キャラクター情報を表す構造体
type Character struct {
	Name     string `json:"name"`
	Rarity   string `json:"rarity"`
	IsPickup bool   `json:"isPickup"`
}

// ガチャの結果を入れる構造体 変数名の先頭が大文字にすると外部からアクセスできる（JSONに変換するために必要）
type GachaResult struct {
	Rarity    string `json:"rarity"`    // レアリティ (`json:"rarity"`は、JSONに変換するときのキー名)
	Character string `json:"character"` // キャラクター名 (`json:"character"`は、JSONに変換するときのキー名)
}

// ブラウザへ返すレスポンス
type GachaResponse struct {
	Results   []GachaResult `json:"results"`   // 今回の結果リスト
	Pity5Star int           `json:"pity5Star"` // 星5天井まであと何回か
	Pity4Star int           `json:"pity4Star"` // 星4天井まであと何回か
	Stones    int           `json:"stones"`    // 所持石数
}
