// 基本レスポンス
export interface ApiResponse {
    success: boolean;
    message: string;
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

// Character構造体
export interface Character {
    id:     number,
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

// Character追加リクエスト
export interface InsertCharacterRequest {
    name:   string,
    rarity: string
}

// 恒常変更のリクエスト型
export interface UpdateConstantRequest {
	charID: number[]
}

// ピックアップ変更のリクエスト型
export interface UpdatePickupRequest {
	bannerID: number,
	star5ID:  number[],
	star4ID:  number[]
}

// 石の付与のリクエスト型
export interface AddStoneRequest {
	uid:    string,
	amount: number,
}