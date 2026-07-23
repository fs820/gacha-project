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
	char_id: number[]
}

// ピックアップ変更のリクエスト型
export interface UpdatePickupRequest {
	banner_id: number,
	star5_id:  number[],
	star4_id:  number[]
}

// 石の付与のリクエスト型
export interface AddStoneRequest {
	uid:    string,
	amount: number,
}