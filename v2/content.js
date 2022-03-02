const baseURL = "https://hiking.biji.co";

function clapAlbum(button) {
  return async function () {
    button.setAttribute("disabled", true);
    const port = chrome.runtime.connect();
    let progress = 0;
    let photoIds = [];

    // listen to progress message
    port.onMessage.addListener(function (message) {
      ({ progress } = message);
      button.innerHTML = `${progress}/${photoIds.length}`;

      // all works are done. close the port
      if (progress === photoIds.length) {
        button.innerHTML = "DONE";
        port.disconnect();
      }
    });

    const albumId = new URLSearchParams(window.location.search).get("album_id");
    // get max page
    const mainPageText = await getAlbumMainPage(albumId);
    const maxPage = parseAlbumMainPage(mainPageText);

    // get all photo IDs
    const pages = Array.from({ length: maxPage }, (_, i) => i + 1);
    photoIds = (
      await Promise.all(pages.map((page) => getAlbumPage(albumId, page)))
    )
      .map((text) => parseAlbumPage(text))
      .flat();

    // send album ID and photo IDs to service worker
    port.postMessage({ albumId, photoIds });

    // change button text to show progress
    button.innerHTML = `0/${photoIds.length}`;
  };
}

async function getAlbumMainPage(albumId) {
  return fetch(
    `https://hiking.biji.co/index.php?q=album&act=photo_list&album_id=${albumId}`
  ).then((resp) => resp.text());
}

async function getAlbumPage(albumId, page) {
  return fetch(
    `https://hiking.biji.co/index.php?q=album&act=photo_list&album_id=${albumId}&page=${page}`
  ).then((resp) => resp.text());
}

// parseAlbumMainPage parses album main max page number
function parseAlbumMainPage(text) {
  const parser = new DOMParser();
  const doc = parser.parseFromString(text, "text/html");

  let maxPage = 1;
  doc.querySelectorAll(".page-item").forEach((item) => {
    const pageLink = item.getAttribute("href");
    if (!pageLink) {
      const p = Number(item.textContent);
      maxPage = Math.max(p, maxPage);
      return;
    }

    const p = new URL(pageLink, baseURL).searchParams.get("page");
    if (p) {
      maxPage = Math.max(p, maxPage);
    }
  });
  return maxPage;
}

// parseAlbumPage parses an album page to get photo IDs
function parseAlbumPage(text) {
  const photoIds = [];

  const parser = new DOMParser();
  const doc = parser.parseFromString(text, "text/html");

  doc.querySelectorAll("li.photo-item").forEach((item) => {
    const photoLink = item.firstElementChild.getAttribute("href");
    if (!photoLink) {
      return;
    }

    const p = new URL(photoLink, baseURL).searchParams.get("ap_id");
    if (p) {
      photoIds.push(p);
    }
  });

  return photoIds;
}

let funcBlock = document.querySelector("div.g-share-wrap");
let button = document.createElement("button");
button.innerHTML = "拍手👏🏽";
button.setAttribute("id", "clap-album");
button.addEventListener("click", clapAlbum(button));
funcBlock.appendChild(button);
