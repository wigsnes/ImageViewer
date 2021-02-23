function deleteFile(file) {
  console.log(file);
  const Http = new XMLHttpRequest();
  const url = "http://localhost:8080/" + file;
  Http.open("DELETE", url);
  Http.send();

  var element = document.getElementById(file);
  element.parentNode.removeChild(element);
}
