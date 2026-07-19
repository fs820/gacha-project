// インターフェースのインポート
import { ApiResponse, GachaBanner, Character, InsertBannerRequest, InsertCharacterRequest, UpdateConstantRequest, UpdatePickupRequest, AddStoneRequest } from "./types.js";

// 最初にチェックボックスを生成
window.onload = async () => {
    await Promise.all([
        createCheckboxes("恒常ガチャ"),
        createBannerInput(),
    ]);

    // バナーの追加ボタン
    addButtonFunc("insertBannerSubmitButton", insertBanner);

    // バナーの変更ボタンがクリックされたときの処理
    addButtonFunc("changeBannerSubmitButton", changeBanner);

    // キャラクター追加ボタンがクリックされたときの処理
    addButtonFunc("insertCharacterSubmitButton", insertCharacter);

    // 恒常の決定ボタンがクリックされたときの処理
    addButtonFunc("constantSubmitButton", changeConstant);

    // ピックアップの決定ボタンがクリックされたときの処理
    addButtonFunc("pickupSubmitButton", changePickUpClick);

    // 石の付与ボタンがクリックされたときの処理
    addButtonFunc("addStoneSubmitButton", addStone);

    // 履歴の削除ボタンがクリックされたときの処理
    addButtonFunc("deleteHistorySubmitButton", deleteHistory);
}

// バナー更新欄を生成する関数
async function createBannerInput() {
    // HTMLから入力を取得する
    const container_gachaBanner = document.getElementById("banner");
    if (!container_gachaBanner) {
        throw new Error("[Error] HTMLにbannerがありません");
    }

    try {
        const banners: GachaBanner[] =
        await postJson<undefined, GachaBanner[]>(`/admin/get_banner`, "POST", undefined);
        banners.forEach(banner => {
            container_gachaBanner.appendChild(createBannerEditor(banner));
        });
    } catch (error) {
        handleError(error);
    }
}

// チェックボックスを生成する関数
async function createCheckboxes(bannerTitle: string) {
    // HTMLから入力を取得する
    const container_star5 = document.getElementById("checkboxContainer_star5");
    const container_star4 = document.getElementById("checkboxContainer_star4");
    if (!(container_star5) || !(container_star4)) {
        throw new Error("[Error] HTMLに[checkboxContainer_star5][checkboxContainer_star4]タグがありません");
    }

    var characters: Character[];
    try {
        characters = await postJson<undefined, Character[]>(`/admin/get_character`, "POST", undefined);
    } catch (error) {
        handleError(error);
        return;
    }

    var pickupIDs: number[];
    try {
        pickupIDs = await postJson<string, number[]>(`/admin/get_pickup_id`, "POST", bannerTitle);
    } catch (error) {
        console.error(error);
        alert("通信エラーが発生しました");
    }

    characters.forEach(char => {
        createCharacterEditor(char, pickupIDs, container_star5, container_star4);
    });
}

// バナーの追加
async function insertBanner() {
    // HTMLから入力を取得する
    const inputTitle = getInput("title");
    const inputCost = getInput("cost");
    const inputProbBaseStar5 = getInput("probBaseStar5");
    const inputProbBaseStar4 = getInput("probBaseStar4");
    const inputStar5Limit = getInput("star5Limit");
    const inputStar4Limit = getInput("star4Limit");
    const inputStar5PickupProb = getInput("star5PickupProb");
    const inputPitySoftStart = getInput("pitySoftStart");
    const inputSoftPityIncrement = getInput("softPityIncrement");

    // 変換する
    const title = inputTitle.value;
    const cost = inputCost.valueAsNumber;
    const probBaseStar5 = inputProbBaseStar5.valueAsNumber;
    const probBaseStar4 = inputProbBaseStar4.valueAsNumber;
    const star5Limit = inputStar5Limit.valueAsNumber;
    const star4Limit = inputStar4Limit.valueAsNumber;
    const star5PickupProb = inputStar5PickupProb.valueAsNumber;
    const pitySoftStart = inputPitySoftStart.valueAsNumber;
    const softPityIncrement = inputSoftPityIncrement.valueAsNumber;
    if (!title || Number.isNaN(cost) || Number.isNaN(probBaseStar5) || Number.isNaN(probBaseStar4) ||
        Number.isNaN(star5Limit) || Number.isNaN(star4Limit) || Number.isNaN(star5PickupProb) ||
        Number.isNaN(pitySoftStart) || Number.isNaN(softPityIncrement)) {
        const apiResponse: ApiResponse = {
            success: false,
            message: "すべてのフィールドを入力してください"
        }
        showResponse(apiResponse);
        return;
    }

    // 渡すデータ
    const requestData: InsertBannerRequest = {
        title: title,
        cost: cost,
        probBaseStar5: probBaseStar5,
        probBaseStar4: probBaseStar4,
        star5Limit: star5Limit,
        star4Limit: star4Limit,
        star5PickupProb: star5PickupProb,
        pitySoftStart: pitySoftStart,
        softPityIncrement: softPityIncrement
    }

    try {
        const apiResponse: ApiResponse =
        await postJson<InsertBannerRequest, ApiResponse>(`/admin/insert_banner`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}

// バナーの更新
async function changeBanner() {
    // HTMLから入力を取得する
    const inputID = getInput("id");
    const inputTitle = getInput("title");
    const inputCost = getInput("cost");
    const inputProbBaseStar5 = getInput("probBaseStar5");
    const inputProbBaseStar4 = getInput("probBaseStar4");
    const inputStar5Limit = getInput("star5Limit");
    const inputStar4Limit = getInput("star4Limit");
    const inputStar5PickupProb = getInput("star5PickupProb");
    const inputPitySoftStart = getInput("pitySoftStart");
    const inputSoftPityIncrement = getInput("softPityIncrement");

    // 変換する
    const id = inputID.valueAsNumber;
    const title = inputTitle.value;
    const cost = inputCost.valueAsNumber;
    const probBaseStar5 = inputProbBaseStar5.valueAsNumber;
    const probBaseStar4 = inputProbBaseStar4.valueAsNumber;
    const star5Limit = inputStar5Limit.valueAsNumber;
    const star4Limit = inputStar4Limit.valueAsNumber;
    const star5PickupProb = inputStar5PickupProb.valueAsNumber;
    const pitySoftStart = inputPitySoftStart.valueAsNumber;
    const softPityIncrement = inputSoftPityIncrement.valueAsNumber;
    if (Number.isNaN(id) || !title || Number.isNaN(cost) || Number.isNaN(probBaseStar5) || Number.isNaN(probBaseStar4) ||
        Number.isNaN(star5Limit) || Number.isNaN(star4Limit) || Number.isNaN(star5PickupProb) ||
        Number.isNaN(pitySoftStart) || Number.isNaN(softPityIncrement)) {
        const apiResponse: ApiResponse = {
            success: false,
            message: "すべてのフィールドを入力してください"
        }
        showResponse(apiResponse);
        return;
    }

    // 渡すデータ
    const requestData: GachaBanner = {
        id: id,
        title: title,
        cost: cost,
        probBaseStar5: probBaseStar5,
        probBaseStar4: probBaseStar4,
        star5Limit: star5Limit,
        star4Limit: star4Limit,
        star5PickupProb: star5PickupProb,
        pitySoftStart: pitySoftStart,
        softPityIncrement: softPityIncrement
    }

    try {
        const apiResponse: ApiResponse =
        await postJson<GachaBanner, ApiResponse>(`/admin/change_banner`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}

// キャラクターの追加
async function insertCharacter() {
    // HTMLから入力を取得する
    const inputName = getInput("charName");
    const inputRarity = getInput("charRarity");

    // 変換する
    const name = inputName.value;
    const rarity = inputRarity.value;
    if (!name || !rarity) {
        const apiResponse: ApiResponse = {
            success: false,
            message: "すべてのフィールドを入力してください"
        }
        showResponse(apiResponse);
        return;
    }

    // 渡すデータ
    const requestData: InsertCharacterRequest = {
        name: name,
        rarity: rarity
    }

    try {
        const apiResponse: ApiResponse =
        await postJson<InsertCharacterRequest, ApiResponse>(`/admin/insert_character`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}

// ピックアップ変更
async function changeConstant() {
    // チェックされているキャラクターを取得
    const selectedCharacter: number[] = [];
    document.querySelectorAll<HTMLInputElement>("#checkboxContainer_constant input[type='checkbox']")
        .forEach(cb => {
            if (cb.checked) {
                selectedCharacter.push(Number(cb.dataset.id));
            }
        });

    // 渡すデータ
    const requestData: UpdateConstantRequest = {
        charID: selectedCharacter
    };

    try {
        const apiResponse: ApiResponse =
        await postJson<UpdateConstantRequest, ApiResponse>(`/admin/update_pickup`, "POST", requestData);
        showResponse(apiResponse);
    } catch (error) {
        handleError(error);
    }
}

// ピックアップ変更
async function changePickUpClick() {
    // バナータイトルを取得する
    const bannerTitle = getInput("bannerTitle");
    const bannerID = Number(bannerTitle.dataset.id)

    // チェックされている星5のキャラクターを取得
    const selectedStar5: number[] = [];
    document.querySelectorAll<HTMLInputElement>("#checkboxContainer_star5 input[type='checkbox']")
        .forEach(cb => {
            if (cb.checked) {
                selectedStar5.push(Number(cb.dataset.id));
            }
        });

    // チェックされている星4のキャラクターを取得
    const selectedStar4: number[] = [];
    document.querySelectorAll<HTMLInputElement>("#checkboxContainer_star4 input[type='checkbox']")
        .forEach(cb => {
            if (cb.checked) {
                selectedStar4.push(Number(cb.dataset.id));
            }
        });

    // ピックアップを変更する
    changePickUp(bannerID, selectedStar5, selectedStar4);
}

// 石の付与
async function addStone() {
    // HTMLから入力を取得する
    const inputUID = getInput("uid");
    const inputAmount = getInput("amount");

    // 変換する
    const uid = inputUID.value;
    const amount = inputAmount.valueAsNumber;
    if (!uid || Number.isNaN(amount)) {
        console.error("すべてのフィールドを入力してください");
        alert("すべてのフィールドを入力してください");
        return;
    }

    // 渡すデータ
    const requestData: AddStoneRequest = {
        uid: uid,
        amount: amount
    };

    try {
        const apiResponse: ApiResponse =
        await postJson<AddStoneRequest, ApiResponse>(`/admin/add_stone`, "POST", requestData);
        showResponse(apiResponse);
    } catch (error) {
        handleError(error);
    }
}

// 履歴の削除
async function deleteHistory() {
    if (!confirm("本当に履歴を削除しますか？")) {
        return;
    }

    try {
        const response: ApiResponse =
        await postJson<undefined, ApiResponse>(`/admin/delete_history`, "POST", undefined);
        showResponse(response);
    } catch (error) {
        handleError(error);
    }
}

// サーバーと通信する関数
async function postJson<TRequest, TResponse>(
    url: string,
    method: string,
    request: TRequest
): Promise<TResponse> {

    const response = await fetch(url, {
        method: method,
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(request)
    });

    if (!response.ok) {
        throw new Error(await response.text());
    }

    return await response.json() as TResponse;
}

// ピックアップ変更
async function changePickUp(bannerID: number, star5ID: number[], star4ID: number[]) {
    // 渡すデータ
    const requestData: UpdatePickupRequest = {
        bannerID: bannerID,
        star5ID: star5ID,
        star4ID: star4ID
    };

    try {
        const apiResponse: ApiResponse =
        await postJson<UpdatePickupRequest, ApiResponse>(`/admin/update_pickup`, "POST", requestData);
        showResponse(apiResponse);
    } catch (error) {
        handleError(error);
    }
}

// キャラクターのエディタを生成する関数
function createCharacterEditor(character: Character, pickupIDs: number[], container_star5: HTMLElement, container_star4: HTMLElement) {
    const label = document.createElement("label");

    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.value = character.name;

    // ピックアップIDにあればチェックしておく
    for (const ID of pickupIDs)
    {
        if (character.id === ID)
        {
            checkbox.checked = true;
        }
    }

    label.appendChild(checkbox);
    label.appendChild(document.createTextNode(character.name));

    if (character.rarity === "星5") {
        container_star5.appendChild(label);
        container_star5.appendChild(document.createElement("br"));
    } else {
        container_star4.appendChild(label);
        container_star4.appendChild(document.createElement("br"));
    }
}

// バナーのエディタを生成する関数
function createBannerEditor(banner: GachaBanner): HTMLElement
{
    const div = document.createElement("div");
    div.className = "banner";

    div.dataset.id = banner.id.toString();
    const idLabel = document.createElement("span");
    idLabel.textContent = `ID: ${banner.id}`;

    div.appendChild(idLabel);
    div.appendChild(createLabeledTextInput("title", "タイトル", banner.title));
    div.appendChild(createLabeledNumberInput("cost", "消費石", banner.cost));
    div.appendChild(createLabeledNumberInput("probBaseStar5", "星5確率", banner.probBaseStar5));
    div.appendChild(createLabeledNumberInput("probBaseStar4", "星4確率", banner.probBaseStar4));
    div.appendChild(createLabeledNumberInput("star5Limit", "星5天井", banner.star5Limit));
    div.appendChild(createLabeledNumberInput("star4Limit", "星4天井", banner.star4Limit));
    div.appendChild(createLabeledNumberInput("star5PickupProb", "星5ピックアップ率", banner.star5PickupProb));
    div.appendChild(createLabeledNumberInput("pitySoftStart", "確率上昇開始回数", banner.pitySoftStart));
    div.appendChild(createLabeledNumberInput("softPityIncrement", "確率上昇率", banner.softPityIncrement));
    return div;
}

// LABELと数値入力欄を生成する関数
function createLabeledNumberInput(name: string, labelText: string, value: number): HTMLElement {

    const wrapper = document.createElement("div");

    const label = document.createElement("label");
    label.textContent = labelText;
    label.htmlFor = name;

    wrapper.appendChild(label);
    wrapper.appendChild(createNumberInput(name, value));

    return wrapper;
}

// LABELとテキスト入力欄を生成する関数
function createLabeledTextInput(name: string, labelText: string, value: string): HTMLElement {

    const wrapper = document.createElement("div");

    const label = document.createElement("label");
    label.textContent = labelText;

    wrapper.appendChild(label);
    wrapper.appendChild(createTextInput(name, value));

    return wrapper;
}

// 数値入力欄を生成する関数
function createNumberInput(name: string, value: number): HTMLInputElement {
    const input = document.createElement("input");
    input.type = "number";
    input.step = "any";
    input.name = name;
    input.id = name;
    input.value = value.toString();
    return input;
}

// テキスト入力欄を生成する関数
function createTextInput(name: string, value: string)
{
    const input = document.createElement("input");
    input.type = "text";
    input.name = name;
    input.value = value;
    return input;
}

// Inputの取得
function getInput(id: string): HTMLInputElement {
    const element = document.getElementById(id);
    if (!(element instanceof HTMLInputElement)) {
        throw new Error(`[Error] HTMLに[id="${id}"]のinput要素がありません`);
    }
    return element;
}

// ボタンに関数を登録する
function addButtonFunc(buttonName:string, func: EventListenerOrEventListenerObject) {
    const button = document.getElementById(buttonName) as HTMLButtonElement;
    if (!(button instanceof HTMLButtonElement))
    {
        throw new Error(`[Error] HTMLに[${buttonName}]タグがありません`);
    }
    button.addEventListener("click", func);
}

// ページ切り替え
function showPage(id: string) {
    // いったん全て非表示
    document.querySelectorAll<HTMLElement>(".page").forEach(page => {
        page.style.display = "none";
    });

    // 指定のページを表示
    const page = document.getElementById(id);
    if (page)
    {
        page.style.display = "block";
    }
}

// レスポンス生成
function showResponse(response: ApiResponse)
{
    console.log(`[${response.success}]${response.message}`);
    alert(`[${response.success}]${response.message}`);
}

// エラーハンドリング
function handleError(error: unknown)
{
    console.error(error);

    const response: ApiResponse = {
        success: false,
        message: "通信エラーが発生しました"
    };

    showResponse(response);
}