const domain = ".biji.co";
let headerCookies;
(async function constructCookieHeader() {
  const cookies = await chrome.cookies.getAll({ domain });
  headerCookies = cookies
    .map((cookie) => `${cookie.name}=${cookie.value}`)
    .join("; ");
})();

// listen to a message from content script
chrome.runtime.onConnect.addListener(function (port) {
  port.onMessage.addListener(async function (message) {
    const { albumId, photoIds } = message;
    // TODO: send back progress
    // port.postMessage({question: "Who's there?"});
    return true;
  });
});

async function clapPhoto(albumId, apId, headerCookies) {
  return fetch("https://hiking.biji.co/album/ajax/clap_photo", {
    headers: {
      accept: "*/*",
      "accept-language": "en-US,en;q=0.9,zh-TW;q=0.8,zh;q=0.7",
      "content-type": "text/plain;charset=UTF-8",
      "sec-ch-ua":
        '"Microsoft Edge";v="95", "Chromium";v="95", ";Not A Brand";v="99"',
      "sec-ch-ua-mobile": "?0",
      "sec-ch-ua-platform": '"macOS"',
      "sec-fetch-dest": "empty",
      "sec-fetch-mode": "cors",
      "sec-fetch-site": "same-origin",
      cookies: headerCookies,
    },
    referrer: `https://hiking.biji.co/index.php?q=album&act=photo&album_id=${albumId}&ap_id=${apId}`,
    referrerPolicy: "no-referrer-when-downgrade",
    body: `{"id":"${apId}"}`,
    method: "POST",
    mode: "cors",
    credentials: "include",
  }).then((resp) => resp.json());
}
