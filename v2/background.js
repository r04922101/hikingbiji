const domain = ".biji.co";

// get cookie to represent the user
let headerCookies;
(async function constructCookieHeader() {
  const cookies = await chrome.cookies.getAll({ domain });
  headerCookies = cookies
    .map((cookie) => `${cookie.name}=${cookie.value}`)
    .join("; ");
})();

// get UA info
function constructUAHeader() {
  const { brands } = navigator.userAgentData;
  return brands
    .map(({ brand, version }) => `"${brand}";v="${version}"`)
    .join(", ");
}
const ua = constructUAHeader();
const { platform } = navigator.userAgentData;

// listen to a message from content script
chrome.runtime.onConnect.addListener(function (port) {
  port.onMessage.addListener(async function (message) {
    const { albumId, photoIds } = message;
    let progress = 0;

    // clap each photo sequentially to make sure each request is handled properly
    for (const photoId of photoIds) {
      await clapPhoto(albumId, photoId, headerCookies);
      progress += 1;
      port.postMessage({ progress });
    }
  });
});

const MAX_RETRY = 3;
async function clapPhoto(albumId, photoId, headerCookies) {
  let ok = false;
  let tryCount = 0;
  // try best efforts to send request with max retrial count
  do {
    try {
      tryCount += 1;
      const { status } = await (
        await fetch("https://hiking.biji.co/album/ajax/clap_photo", {
          headers: {
            accept: "*/*",
            "accept-language": "en-US,en;q=0.9,zh-TW;q=0.8,zh;q=0.7",
            "content-type": "text/plain;charset=UTF-8",
            "sec-ch-ua": ua,
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": platform,
            "sec-fetch-dest": "empty",
            "sec-fetch-mode": "cors",
            "sec-fetch-site": "same-origin",
            cookies: headerCookies,
          },
          referrer: `https://hiking.biji.co/index.php?q=album&act=photo&album_id=${albumId}&ap_id=${photoId}`,
          referrerPolicy: "no-referrer-when-downgrade",
          body: `{"id":"${photoId}"}`,
          method: "POST",
          mode: "cors",
          credentials: "include",
        })
      ).json();
      ok = true;
      console.log(`finished clapping photo ${photoId} with status ${status}`);
    } catch (err) {
      console.error(
        `failed to clap photo ${photoId} with trial ${tryCount} times: ${err}`
      );
    }
  } while (!ok && tryCount < MAX_RETRY);
}
