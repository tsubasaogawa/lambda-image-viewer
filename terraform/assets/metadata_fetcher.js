document.addEventListener("DOMContentLoaded", async function () {
    try {
        // <img> タグを取得
        const img = document.getElementById("targetImage");
        if (img) {
            // 画像ファイル名を取得
            const src = img.getAttribute("src");
            const fileName = "blog/IMG_20240508_130440.jpg"; // src.split('/').pop();

            // API エンドポイント
            const apiEndpoint = "https://photoview.ogatube.com/metadata";

            // API へのリクエスト
            const response = await fetch(`${apiEndpoint}/${fileName}`);
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            const data = await response.json();

            // レスポンスデータを表示する要素を作成
            // const metadataDiv = document.createElement("div");
            // metadataDiv.innerText = `Metadata: ${JSON.stringify(data)}`;

            // <img> タグのすぐ下に挿入
            // img.parentNode.insertBefore(metadataDiv, img.nextSibling);
            console.log(JSON.stringify(data));
        }
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
});