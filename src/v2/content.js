function clapAlbum() {
  alert("clapping album");
}

let sns = document.querySelector("div.sns-block");
let button = document.createElement("button");
button.innerHTML = "Clap album";
button.setAttribute("id", "clap-album");
button.addEventListener("click", clapAlbum);
sns.appendChild(button);
