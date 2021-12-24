// listen to a message from content scriipt
chrome.runtime.onConnect.addListener(function(port) {
  port.onMessage.addListener(function(message) {
    const { start } = message;
    if (start) {
      // start clapping album
    }
    // TODO: send back progress
    // port.postMessage({question: "Who's there?"});
  });
});