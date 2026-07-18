// 最初にチェックボックスを生成
window.onload = async function () {
    createCheckboxes();
};
// キャラクターの追加
async function insertCharacter() {
    // HTMLから入力を取得する
    const inputName = document.getElementById("charName");
    const inputRarity = document.getElementById("charRarity");
    if (!(inputName instanceof HTMLInputElement) || !(inputRarity instanceof HTMLInputElement)) {
        alert("[Error] HTMLにタグがありません");
        return;
    }
    // 変換する
    const name = inputName.value;
    const rarity = inputRarity.value;
    if (!name || !rarity) {
        alert("すべてのフィールドを入力してください");
        return;
    }
    // 渡すデータ
    const requestData = {
        name: name,
        rarity: rarity
    };
    try {
        const response = await fetch(`/admin/insert_character`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });
        const text = await response.text();
        alert(text);
    }
    catch (error) {
        alert("通信エラーが発生しました");
    }
}
// ピックアップ変更
async function changePickUp(bannerTitle, rarity, names) {
    // 渡すデータ
    const requestData = {
        bannerTitle: bannerTitle,
        rarity: rarity,
        names: names
    };
    try {
        const response = await fetch(`/admin/update_pickup`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });
        const text = await response.text();
        alert(text);
    }
    catch (error) {
        alert("通信エラーが発生しました");
    }
}
// ピックアップの決定ボタンがクリックされたときの処理
document.getElementById("submitButton").addEventListener("click", () => {
    // チェックされている星5のキャラクターを取得
    const selected = [];
    document.querySelectorAll("#checkboxContainer_star5 input[type='checkbox']")
        .forEach(cb => {
        if (cb.checked) {
            selected.push(cb.value);
        }
    });
    changePickUp("星5", selected.join(","));
    // チェックされている星4のキャラクターを取得
    const selected2 = [];
    document.querySelectorAll("#checkboxContainer_star4 input[type='checkbox']")
        .forEach(cb => {
        if (cb.checked) {
            selected2.push(cb.value);
        }
    });
    changePickUp("星4", selected2.join(","));
});
// 石の付与
async function addStone() {
    const uid = document.getElementById("uid").value;
    const amount = document.getElementById("amount").value;
    if (!uid || !amount) {
        alert("すべてのフィールドを入力してください");
        return;
    }
    // 渡すデータ
    const requestData = {
        uid: uid,
        amount: amount
    };
    try {
        const response = await fetch(`/admin/add_stone`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });
        const text = await response.text();
        alert(text);
    }
    catch (error) {
        alert("通信エラーが発生しました");
    }
}
// 履歴の削除
async function deleteHistory() {
    if (!confirm("本当に履歴を削除しますか？")) {
        return;
    }
    try {
        const response = await fetch(`/admin/delete_history`, {
            method: "POST",
        });
        const text = await response.text();
        alert(text);
    }
    catch (error) {
        alert("通信エラーが発生しました");
    }
}
// ページ切り替え
function showPage(id) {
    document.querySelectorAll(".page").forEach(page => {
        page.style.display = "none";
    });
    document.getElementById(id).style.display = "block";
}
// チェックボックスを生成する関数
async function createCheckboxes() {
    const container_star5 = document.getElementById("checkboxContainer_star5");
    const container_star4 = document.getElementById("checkboxContainer_star4");
    items_star5 = [];
    items_star4 = [];
    try {
        const response = await fetch(`/admin/get_character`, {
            method: "POST",
        });
        const data = await response.json();
        for (const item of data) {
            const newItem = {
                id: item.id,
                name: item.name,
                rarity: item.rarity,
            };
            if (newItem.rarity === "星5") {
                items_star5.push(newItem);
            }
            else if (newItem.rarity === "星4") {
                items_star4.push(newItem);
            }
        }
    }
    catch (error) {
        alert("通信エラーが発生しました");
    }
    items_star5.forEach(item => {
        const label = document.createElement("label");
        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.value = item.name;
        checkbox.checked = item.isPickup; // isPickupがtrueの場合はチェックを入れる
        label.appendChild(checkbox);
        label.appendChild(document.createTextNode(item.name));
        container_star5.appendChild(label);
        container_star5.appendChild(document.createElement("br"));
    });
    items_star4.forEach(item => {
        const label = document.createElement("label");
        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.value = item.name;
        checkbox.checked = item.isPickup; // isPickupがtrueの場合はチェックを入れる
        label.appendChild(checkbox);
        label.appendChild(document.createTextNode(item.name));
        container_star4.appendChild(label);
        container_star4.appendChild(document.createElement("br"));
    });
}
export {};
//# sourceMappingURL=script.js.map