document.addEventListener("DOMContentLoaded", async function () {
    if (typeof mfApiEndpoint === 'undefined') {
        throw new Error('mfApiEndpoint is required');
    }

    const imgs = document.getElementsByClassName("load-exif");
    const getPathFromUrl = function (url) {
        return url.split("/").slice(3).join("/");
    };
    const formatMetadata = function (data) {
        return `${data["camera"]} (${data["lens"]}), ${data["focal_length"]}mm, f/${data["f"]}, ISO ${data["iso"]}, ${data["shutter_speed"]} sec, ${data["exposure"]} EV`;
    };

    const requests = Array.from(imgs).map(async (img) => {
        const src = img.getAttribute("src");

        const response = await fetch(`${mfApiEndpoint}/${getPathFromUrl(src)}`);
        if (!response.ok) {
            console.error('Network response was not ok: ' + response.statusText);
        }
        const metadata = await response.json();

        const metadataDiv = document.createElement("div");
        metadataDiv.className = 'exif';
        metadataDiv.innerHTML = `<i>${formatMetadata(metadata)}</i>`;

        img.parentNode.insertBefore(metadataDiv, img.nextSibling);
    });

    await Promise.all(requests);
});
