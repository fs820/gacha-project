package core

// ユーザデータ
type UserData struct {
	UID                    string        `json:"uid"`
	Stones                 int           `json:"stones"`
	Star4LimitCounter      int           `json:"star4LimitCounter"`
	Star5LimitCounter      int           `json:"star5LimitCounter"`
	IsNextPickupGuaranteed bool          `json:"isNextPickupGuaranteed"`
	GachaHistory           []GachaResult `json:"gachaHistory"`
}

// キャラクター情報を表す構造体
type Character struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

// ガチャの結果を入れる構造体 変数名の先頭が大文字にすると外部からアクセスできる（JSONに変換するために必要）
type GachaResult struct {
	Character Character `json:"character"` // キャラクターの全情報
}

// ブラウザへ返すレスポンス
type GachaResponse struct {
	Results   []GachaResult `json:"results"`   // 今回の結果リスト
	Pity5Star int           `json:"pity5Star"` // 星5天井まであと何回か
	Pity4Star int           `json:"pity4Star"` // 星4天井まであと何回か
	Stones    int           `json:"stones"`    // 所持石数
}

// ガチャの基本データ
type GachaBanner struct {
	ID                int    `json:"id"`
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

// 現在ピックアップ中のキャラ一覧
type PickupCharacters struct {
	Star5 []Character `json:"star5"`
	Star4 []Character `json:"star4"`
}
