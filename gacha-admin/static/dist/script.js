// 最初にチェックボックスを生成
window.onload = async () => {
    // バナーとキャラクターの取得
    let banners;
    let characters;
    try {
        [banners, characters] = await Promise.all([
            postJson(`/admin/get_banner`, "POST", undefined),
            postJson(`/admin/get_character`, "POST", undefined),
        ]);
    }
    catch (error) {
        handleError(error);
        return; // banners/charactersが無いと以降の初期化ができないのでここで止める
    }
    await Promise.all([
        createBannerInput(banners),
        createConstantCheckboxes(characters),
        createPickupCheckboxes(banners, characters),
    ]);
    // タブ切り替えボタン
    addShowPageEvent();
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
    // バナーのプルダウンが変更されたときの処理
    addpulldownFunc("change", "banner_select", updateChangeBanner);
    addpulldownFunc("pulldownContainer_banner", "banner_select", updatePickupCheckboxes);
};
// バナー更新欄を生成する関数
async function createBannerInput(banners) {
    // HTMLから入力を取得する
    const container_gachaBanner = document.getElementById("change");
    if (!container_gachaBanner) {
        throw new Error("[Error] HTMLにbannerがありません");
    }
    // プルダウンを作る
    createBannerPulldown("change", banners);
    if (!banners || banners.length <= 0) {
        return;
    }
    container_gachaBanner.appendChild(createBannerEditor("change_banner", banners[0]));
}
// チェックボックスを生成する関数
async function createConstantCheckboxes(characters) {
    // HTMLから入力を取得する
    const container_constant = document.getElementById("checkboxContainer_constant");
    if (!(container_constant)) {
        throw new Error("[Error] HTMLに[container_constant]タグがありません");
    }
    var nowIDs;
    try {
        nowIDs = await postJson(`/admin/get_constant_id`, "POST", undefined);
    }
    catch (error) {
        handleError(error);
        return;
    }
    characters.forEach(char => {
        createConstantCharacterEditor(char, nowIDs, container_constant);
    });
}
// チェックボックスを生成する関数
async function createPickupCheckboxes(banners, characters) {
    // バナープルダウンを作る
    createBannerPulldown("pulldownContainer_banner", banners);
    // HTMLから入力を取得する
    const container_star5 = document.getElementById("checkboxContainer_star5");
    const container_star4 = document.getElementById("checkboxContainer_star4");
    if (!(container_star5) || !(container_star4)) {
        throw new Error("[Error] HTMLに[checkboxContainer_star5][checkboxContainer_star4]タグがありません");
    }
    var pickupIDs;
    try {
        pickupIDs = await postJson(`/admin/get_pickup_id`, "POST", 0);
    }
    catch (error) {
        handleError(error);
        return;
    }
    characters.forEach(char => {
        createPickupCharacterEditor(char, pickupIDs, container_star5, container_star4);
    });
}
// バナーの追加
async function insertBanner() {
    // HTMLから入力を取得する
    const inputTitle = getFormToInput("insert_banner", "title");
    const inputCost = getFormToInput("insert_banner", "cost");
    const inputProbBaseStar5 = getFormToInput("insert_banner", "probBaseStar5");
    const inputProbBaseStar4 = getFormToInput("insert_banner", "probBaseStar4");
    const inputStar5Limit = getFormToInput("insert_banner", "star5Limit");
    const inputStar4Limit = getFormToInput("insert_banner", "star4Limit");
    const inputStar5PickupProb = getFormToInput("insert_banner", "star5PickupProb");
    const inputPitySoftStart = getFormToInput("insert_banner", "pitySoftStart");
    const inputSoftPityIncrement = getFormToInput("insert_banner", "softPityIncrement");
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
        const apiResponse = {
            success: false,
            message: "すべてのフィールドを入力してください"
        };
        showResponse(apiResponse);
        return;
    }
    // 渡すデータ
    const requestData = {
        title: title,
        cost: cost,
        probBaseStar5: probBaseStar5,
        probBaseStar4: probBaseStar4,
        star5Limit: star5Limit,
        star4Limit: star4Limit,
        star5PickupProb: star5PickupProb,
        pitySoftStart: pitySoftStart,
        softPityIncrement: softPityIncrement
    };
    try {
        const apiResponse = await postJson(`/admin/insert_banner`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}
// バナーの更新
async function changeBanner() {
    const change_banner = document.getElementById("change_banner");
    if (!(change_banner)) {
        throw new Error("[Error] HTMLに[change_banner]タグがありません");
    }
    const banner_select = change_banner.querySelector(`[name='banner_select']`);
    if (!banner_select) {
        throw new Error("[Error] change_banner[banner_select]タグがありません");
    }
    const inputTitle = getFormToInput("change_banner", "title");
    const inputCost = getFormToInput("change_banner", "cost");
    const inputProbBaseStar5 = getFormToInput("change_banner", "probBaseStar5");
    const inputProbBaseStar4 = getFormToInput("change_banner", "probBaseStar4");
    const inputStar5Limit = getFormToInput("change_banner", "star5Limit");
    const inputStar4Limit = getFormToInput("change_banner", "star4Limit");
    const inputStar5PickupProb = getFormToInput("change_banner", "star5PickupProb");
    const inputPitySoftStart = getFormToInput("change_banner", "pitySoftStart");
    const inputSoftPityIncrement = getFormToInput("change_banner", "softPityIncrement");
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
        const apiResponse = {
            success: false,
            message: "すべてのフィールドを入力してください"
        };
        showResponse(apiResponse);
        return;
    }
    // 渡すデータ
    const requestData = {
        id: Number(banner_select.value),
        title: title,
        cost: cost,
        probBaseStar5: probBaseStar5,
        probBaseStar4: probBaseStar4,
        star5Limit: star5Limit,
        star4Limit: star4Limit,
        star5PickupProb: star5PickupProb,
        pitySoftStart: pitySoftStart,
        softPityIncrement: softPityIncrement
    };
    try {
        const apiResponse = await postJson(`/admin/change_banner`, "POST", requestData);
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
    const inputRarity = getSelect("charRarity");
    // 変換する
    const name = inputName.value;
    const rarity = inputRarity.value;
    if (!name || !rarity) {
        const apiResponse = {
            success: false,
            message: "すべてのフィールドを入力してください"
        };
        showResponse(apiResponse);
        return;
    }
    // 渡すデータ
    const requestData = {
        name: name,
        rarity: rarity
    };
    try {
        const apiResponse = await postJson(`/admin/insert_character`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}
// ピックアップ変更
async function changeConstant() {
    // チェックされているキャラクターを取得
    const selectedCharacter = [];
    document.querySelectorAll("#checkboxContainer_constant input[type='checkbox']")
        .forEach(cb => {
        if (cb.checked) {
            selectedCharacter.push(Number(cb.dataset.id));
        }
    });
    // 渡すデータ
    const requestData = {
        charID: selectedCharacter
    };
    try {
        const apiResponse = await postJson(`/admin/update_constant`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}
// ピックアップ変更
async function changePickUpClick() {
    const pulldownContainer_banner = document.getElementById("pulldownContainer_banner");
    if (!(pulldownContainer_banner)) {
        throw new Error("[Error] HTMLに[pulldownContainer_banner]タグがありません");
    }
    const banner_select = pulldownContainer_banner.querySelector(`[name='banner_select']`);
    if (!banner_select) {
        throw new Error("[Error] pulldownContainer_bannerに[banner_select]タグがありません");
    }
    // チェックされている星5のキャラクターを取得
    const selectedStar5 = [];
    document.querySelectorAll("#checkboxContainer_star5 input[type='checkbox']")
        .forEach(cb => {
        if (cb.checked) {
            selectedStar5.push(Number(cb.dataset.id));
        }
    });
    // チェックされている星4のキャラクターを取得
    const selectedStar4 = [];
    document.querySelectorAll("#checkboxContainer_star4 input[type='checkbox']")
        .forEach(cb => {
        if (cb.checked) {
            selectedStar4.push(Number(cb.dataset.id));
        }
    });
    // ピックアップを変更する
    changePickUp(Number(banner_select.value), selectedStar5, selectedStar4);
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
    const requestData = {
        uid: uid,
        amount: amount
    };
    try {
        const apiResponse = await postJson(`/admin/add_stone`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}
// 履歴の削除
async function deleteHistory() {
    if (!confirm("本当に履歴を削除しますか？")) {
        return;
    }
    try {
        const response = await postJson(`/admin/delete_history`, "POST", undefined);
        showResponse(response);
    }
    catch (error) {
        handleError(error);
    }
}
// サーバーと通信する関数
async function postJson(url, method, request) {
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
    return await response.json();
}
// ピックアップ変更
async function changePickUp(bannerID, star5ID, star4ID) {
    // 渡すデータ
    const requestData = {
        bannerID: bannerID,
        star5ID: star5ID,
        star4ID: star4ID
    };
    try {
        const apiResponse = await postJson(`/admin/update_pickup`, "POST", requestData);
        showResponse(apiResponse);
    }
    catch (error) {
        handleError(error);
    }
}
// 恒常キャラクターのエディタを生成する関数
function createConstantCharacterEditor(character, nowIDs, container_constant) {
    const label = document.createElement("label");
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.dataset.id = character.id.toString();
    checkbox.value = character.name;
    // 今のIDにあればチェックしておく
    for (const ID of nowIDs) {
        if (character.id === ID) {
            checkbox.checked = true;
        }
    }
    label.appendChild(checkbox);
    label.appendChild(document.createTextNode(character.name));
    container_constant.appendChild(label);
    container_constant.appendChild(document.createElement("br"));
}
// ピックアップキャラクターのエディタを生成する関数
function createPickupCharacterEditor(character, pickupIDs, container_star5, container_star4) {
    const label = document.createElement("label");
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.dataset.id = character.id.toString();
    checkbox.value = character.name;
    // ピックアップIDにあればチェックしておく
    checkbox.checked = pickupIDs.includes(character.id);
    label.appendChild(checkbox);
    label.appendChild(document.createTextNode(character.name));
    if (character.rarity === "星5") {
        container_star5.appendChild(label);
        container_star5.appendChild(document.createElement("br"));
    }
    else {
        container_star4.appendChild(label);
        container_star4.appendChild(document.createElement("br"));
    }
}
// バナーのエディタを生成する関数
function createBannerEditor(formName, banner) {
    const form = document.getElementById(formName);
    if (!(form instanceof HTMLFormElement)) {
        throw new Error(`[Error] HTMLに[${formName}]タグがありません`);
    }
    form.replaceChildren();
    form.id = "change_banner";
    form.dataset.id = banner.id.toString();
    const idLabel = document.createElement("span");
    idLabel.textContent = `ID: ${banner.id}`;
    form.appendChild(idLabel);
    form.appendChild(createLabeledTextInput("title", "タイトル", banner.title));
    form.appendChild(createLabeledNumberInput("cost", "消費石", banner.cost));
    form.appendChild(createLabeledNumberInput("probBaseStar5", "星5確率", banner.probBaseStar5));
    form.appendChild(createLabeledNumberInput("probBaseStar4", "星4確率", banner.probBaseStar4));
    form.appendChild(createLabeledNumberInput("star5Limit", "星5天井", banner.star5Limit));
    form.appendChild(createLabeledNumberInput("star4Limit", "星4天井", banner.star4Limit));
    form.appendChild(createLabeledNumberInput("star5PickupProb", "星5ピックアップ率", banner.star5PickupProb));
    form.appendChild(createLabeledNumberInput("pitySoftStart", "確率上昇開始回数", banner.pitySoftStart));
    form.appendChild(createLabeledNumberInput("softPityIncrement", "確率上昇率", banner.softPityIncrement));
    return form;
}
// LABELと数値入力欄を生成する関数
function createLabeledNumberInput(name, labelText, value) {
    const wrapper = document.createElement("div");
    const label = document.createElement("label");
    label.textContent = labelText;
    label.htmlFor = name;
    wrapper.appendChild(label);
    wrapper.appendChild(createNumberInput(name, value));
    return wrapper;
}
// LABELとテキスト入力欄を生成する関数
function createLabeledTextInput(name, labelText, value) {
    const wrapper = document.createElement("div");
    const label = document.createElement("label");
    label.textContent = labelText;
    wrapper.appendChild(label);
    wrapper.appendChild(createTextInput(name, value));
    return wrapper;
}
// 数値入力欄を生成する関数
function createNumberInput(name, value) {
    const input = document.createElement("input");
    input.type = "number";
    input.step = "any";
    input.name = name;
    input.value = value.toString();
    return input;
}
// テキスト入力欄を生成する関数
function createTextInput(name, value) {
    const input = document.createElement("input");
    input.type = "text";
    input.name = name;
    input.value = value;
    return input;
}
// タブ切り替えボタンに切り替え機能を付与する
function addShowPageEvent() {
    const tabs = document.querySelector(".tabs");
    if (!tabs) {
        throw new Error("[Error] HTMLにtabsがありません");
    }
    const buttons = tabs.querySelectorAll("button");
    buttons.forEach(button => {
        button.addEventListener("click", () => {
            showPage(button.name);
        });
    });
}
// バナープルダウンの生成
function createBannerPulldown(containerName, banners) {
    const container = document.getElementById(containerName);
    if (!(container)) {
        throw new Error(`[Error] HTMLに[${containerName}]タグがありません`);
    }
    if (banners.length === 0) {
        alert("バナーがありません！先にバナーを追加してください");
        return;
    }
    const banner_select = document.createElement("select");
    banner_select.name = "banner_select";
    banners.forEach(banner => {
        const option = document.createElement("option");
        option.value = banner.id.toString(); // 送信する値
        option.textContent = banner.title; // 表示する文字
        banner_select.appendChild(option);
    });
    container.appendChild(banner_select);
}
// バナー更新欄の更新
async function updateChangeBanner() {
    // HTMLから入力を取得する
    const container_gachaBanner = document.getElementById("change");
    if (!container_gachaBanner) {
        throw new Error("[Error] HTMLにbannerがありません");
    }
    const pulldown = container_gachaBanner.querySelector(`[name="banner_select"]`);
    if (!pulldown) {
        throw new Error(`[Error] change_bannerに[banner_select]タグがありません`);
    }
    const id = Number(pulldown.value);
    try {
        const banners = await postJson(`/admin/get_banner`, "POST", undefined);
        const banner = banners.find(b => b.id === id);
        if (!banner) {
            throw new Error(`[Error] バナー更新エラー`);
        }
        createBannerEditor("change_banner", banner);
    }
    catch (error) {
        handleError(error);
    }
}
// ピックアップチェックボックスの更新
async function updatePickupCheckboxes() {
    const pulldownContainer_banner = document.getElementById("pulldownContainer_banner");
    if (!(pulldownContainer_banner)) {
        throw new Error("[Error] HTMLに[pulldownContainer_banner]タグがありません");
    }
    const banner_select = pulldownContainer_banner.querySelector(`[name='banner_select']`);
    if (!banner_select) {
        throw new Error("[Error] pulldownContainer_bannerに[banner_select]タグがありません");
    }
    // HTMLから範囲を取得する
    const container_star5 = document.getElementById("checkboxContainer_star5");
    const container_star4 = document.getElementById("checkboxContainer_star4");
    if (!(container_star5) || !(container_star4)) {
        throw new Error("[Error] HTMLに[checkboxContainer_star5][checkboxContainer_star4]タグがありません");
    }
    var pickupIDs;
    try {
        pickupIDs = await postJson(`/admin/get_pickup_id`, "POST", Number(banner_select.value));
    }
    catch (error) {
        handleError(error);
        return;
    }
    const checkboxes_star5 = container_star5.querySelectorAll("input[type='checkbox']");
    checkboxes_star5.forEach(cb => {
        const id = Number(cb.dataset.id);
        cb.checked = pickupIDs.includes(id);
    });
    const checkboxes_star4 = container_star4.querySelectorAll("input[type='checkbox']");
    checkboxes_star4.forEach(cb => {
        const id = Number(cb.dataset.id);
        cb.checked = pickupIDs.includes(id);
    });
}
// Inputの取得
function getInput(id) {
    const element = document.getElementById(id);
    if (!(element instanceof HTMLInputElement)) {
        throw new Error(`[Error] HTMLに[id="${id}"]のinput要素がありません`);
    }
    return element;
}
// Inputの取得
function getSelect(id) {
    const element = document.getElementById(id);
    if (!(element instanceof HTMLSelectElement)) {
        throw new Error(`[Error] HTMLに[id="${id}"]のselect要素がありません`);
    }
    return element;
}
// Inputの取得
function getFormToInput(formName, name) {
    const form = document.getElementById(formName);
    if (!(form)) {
        throw new Error(`[Error] HTMLにform:${formName}要素がありません`);
    }
    const element = form.querySelector(`[name="${name}"]`);
    if (!(element instanceof HTMLInputElement)) {
        throw new Error(`[Error] HTMLのform:${formName}に[name="${name}"]のinput要素がありません`);
    }
    return element;
}
// ボタンに関数を登録する
function addButtonFunc(buttonName, func) {
    const button = document.getElementById(buttonName);
    if (!(button instanceof HTMLButtonElement)) {
        throw new Error(`[Error] HTMLに[${buttonName}]タグがありません`);
    }
    button.addEventListener("click", func);
}
// ボタンに関数を登録する
function addpulldownFunc(containerName, pulldownName, func) {
    const container = document.getElementById(containerName);
    if (!container) {
        throw new Error(`[Error] HTMLに[${containerName}]タグがありません`);
    }
    const pulldown = container.querySelector(`[name="${pulldownName}"]`);
    if (!pulldown) {
        throw new Error(`[Error] ${containerName}に[${pulldownName}]タグがありません`);
    }
    pulldown.addEventListener("change", func);
}
// ページ切り替え
function showPage(id) {
    // いったん全て非表示
    document.querySelectorAll(".page").forEach(page => {
        page.style.display = "none";
    });
    // 指定のページを表示
    const page = document.getElementById(id);
    if (page) {
        page.style.display = "block";
    }
}
// レスポンス生成
function showResponse(response) {
    console.log(`[${response.success}]${response.message}`);
    alert(`[${response.success}]${response.message}`);
}
// エラーハンドリング
function handleError(error) {
    console.error(error);
    const response = {
        success: false,
        message: "通信エラーが発生しました"
    };
    showResponse(response);
}
export {};
//# sourceMappingURL=script.js.map