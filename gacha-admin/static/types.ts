// Character構造体
export interface Character {
    id:     number,
    name:   string,
    rarity: string
}

// ガチャバナー構造体
export interface GachaBanner {
	id:                number,
	title:             string,
	cost:              number,
	probBaseStar5:     number,
	probBaseStar4:     number,
	star5Limit:        number,
	star4Limit:        number,
	star5PickupProb:   number,
	pitySoftStart:     number,
	softPityIncrement: number,
}

// 基本レスポンス
export interface ApiResponse {
    success: boolean;
    message: string;
}

// Character追加リクエスト
export interface InsertCharacterRequest {
    name:   string,
    rarity: string
}

// ガチャの基本データ
export interface InsertBannerRequest {
	title:             string,
	cost:              number,
	probBaseStar5:     number,
	probBaseStar4:     number,
	star5Limit:        number,
	star4Limit:        number,
	star5PickupProb:   number,
	pitySoftStart:     number,
	softPityIncrement: number
}

// ピックアップ変更のリクエスト型
export interface UpdatePickupRequest {
	bannerTitle: string,
	star5Names:  string[],
	star4Names:  string[]
}

// 石の付与のリクエスト型
export interface AddStoneRequest {
	uid:    string,
	amount: number,
}