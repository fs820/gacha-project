// 最初にチェックボックスを生成
window.onload = async function() {
    const container_star5 = document.getElementById("checkboxContainer_star5");
    const container_star4 = document.getElementById("checkboxContainer_star4");

    items_star5 = [];
    items_star4 = [];

    try {
        const response = await fetch(`/admin/get_character`, {
            method: "POST",
            headers: {
                Authorization: `Bearer supersecret`,
            }
        });
        const data = await response.json();

        for (const item of data) {
            const newItem = {
            name: item.name,
            rarity: item.rarity,
            isPickup: item.isPickup
            }
            if (newItem.rarity === "星5") {
                items_star5.push(newItem);
            } else if (newItem.rarity === "星4") {
                items_star4.push(newItem);
            }
        }
    } catch (error) {
        alert("通信エラーが発生しました");
    }

    items_star5.forEach(item => {
        const label = document.createElement("label");

        const checkbox = document.createElement("input");
        checkbox.type = "radio";
        checkbox.name = "star5";
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

    // チェックボックスの制限を設定
    setupCheckboxLimit()
}

// ページ切り替え
function showPage(id) {
    document.querySelectorAll(".page").forEach(page => {
        page.style.display = "none";
    });

    document.getElementById(id).style.display = "block";
}

// チェックボックスの制限を設定
function setupCheckboxLimit() {
    // 星4
    document.querySelectorAll("#checkboxContainer_star4 input[type='checkbox']")
        .forEach(cb => {
            cb.addEventListener("change", () => {
                const checked = document.querySelectorAll(
                    "#checkboxContainer_star4 input[type='checkbox']:checked"
                );

                if (checked.length > 3) {
                    cb.checked = false;
                    alert("星4は3体まで選択できます。");
                }
            });
        });
}

// 決定ボタンがクリックされたときの処理
document.getElementById("submitButton").addEventListener("click", () => {
    // チェックされている星5のキャラクター1体を取得
    const selected = [];
    document.querySelectorAll("#checkboxContainer_star5 input[type='radio']")
        .forEach(cb => {
            if (cb.checked) {
                selected.push(cb.value);
            }
        });
    if (selected.length == 1) {
        changePickUp("星5", selected[0]);
    } else {
        alert("星5のキャラクターは1体だけ選択してください");
        return;
    }

    // チェックされている星4のキャラクター3体までを取得
    const selected2 = [];
    document.querySelectorAll("#checkboxContainer_star4 input[type='checkbox']")
        .forEach(cb => {
            if (cb.checked) {
                selected2.push(cb.value);
            }
        });
    if (selected2.length > 0 && selected2.length <= 3) {
        changePickUp("星4", selected2.join(","));
    } else {
        alert("星4のキャラクターは1体以上3体まで選択してください");
        return;
    }
});

async function insertCharacter() {
    const name = document.getElementById("charName").value;
    const rarity = document.getElementById("charRarity").value;

    if (!name || !rarity) {
        alert("すべてのフィールドを入力してください");
        return;
    }

    try {
        const response = await fetch(`/admin/insert_character?name=${encodeURIComponent(name)}&rarity=${encodeURIComponent(rarity)}`, {
            method: "POST",
            headers: {
                Authorization: `Bearer supersecret`,
            }
        });
        const text = await response.text();
        alert(text);
    }
    catch (error) {
        alert("通信エラーが発生しました");
    }
}

// ピックアップ変更
async function changePickUp(rarity, names) {
    try {
        const response = await fetch(`/admin/update_pickup?rarity=${encodeURIComponent(rarity)}&name=${encodeURIComponent(names)}`, {
            method: "POST",
            headers: {
                Authorization: `Bearer supersecret`,
            }
        });
        const text = await response.text();
        alert(text);
    } catch (error) {
        alert("通信エラーが発生しました");
    }
}

// 石の付与
async function addStone() {
    const uid = document.getElementById("uid").value;
    const amount = document.getElementById("amount").value;

    if (!uid || !amount) {
        alert("すべてのフィールドを入力してください");
        return;
    }

    try {
        const response = await fetch(`/admin/add_stone?uid=${encodeURIComponent(uid)}&amount=${encodeURIComponent(amount)}`, {
            method: "POST",
            headers: {
                Authorization: `Bearer supersecret`,
            }
        });
        const text = await response.text();
        alert(text);
    } catch (error) {
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
            headers: {
                Authorization: `Bearer supersecret`,
            }
        });
        const text = await response.text();
        alert(text);
    } catch (error) {
        alert("通信エラーが発生しました");
    }
}