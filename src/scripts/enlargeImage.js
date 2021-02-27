function removeImage() {
  var imageWrapper = document.getElementById("enlargedImageWrapper");
  imageWrapper.remove();
}

function enlargeImage(event) {
  var div = document.createElement("div");
  var img = document.createElement("img");
  const page = document.getElementById("page");

  div.id = "enlargedImageWrapper";
  div.onclick = removeImage;

  img.src = event.src;
  img.id = "enlargedImage"; // height: auto; width: auto; position: fixed; left: 39%; top: 0;

  div.appendChild(img);
  page.appendChild(div);
}
