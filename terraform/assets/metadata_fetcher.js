document.addEventListener("DOMContentLoaded", async function () {
    const apiEndpoint = "https://photoview.ogatube.com/metadata";
    const imgs = document.getElementsByClassName("load-exif");
    const getPathFromUrl = function (url) {
        return url.split("/").slice(3).join("/");
    };
    const formatMetadata = function (data) {
        return `${data["camera"]} (${data["lens"]}), ${data["focal_length"]}mm, f/${data["f"]}, ISO ${data["iso"]}, ${data["shutter_speed"]} sec, ${data["exposure"]} EV`;
    };

    for (let img of imgs) {
        const src = img.getAttribute("src");

        const response = await fetch(`${apiEndpoint}/${getPathFromUrl(src)}`);
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        const metadata = await response.json();

        const metadataDiv = document.createElement("div");
        metadataDiv.innerHTML = `<i>${formatMetadata(metadata)}</i>`;

        img.parentNode.insertBefore(metadataDiv, img.nextSibling);
    }
});