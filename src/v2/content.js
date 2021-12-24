
function clapAlbum() {
  alert("clapping album");
  const port = chrome.runtime.connect({ name: "hikingbiji" });

  // TODO: handle message from service worker
  // TODO: show progress on web
  port.onMessage.addListener(function (message) {
  });

  // TODO: handle disconnect event
  port.onDisconnect.addListener(() => {
    alert("connection disconnect");
  });
  
  // post message to instruct to start clapping the album
  port.postMessage({ start: true });
}

let sns = document.querySelector("div.sns-block");
let button = document.createElement("button");
button.innerHTML = "Clap album";
button.setAttribute("id", "clap-album");
button.addEventListener("click", clapAlbum);
sns.appendChild(button);
