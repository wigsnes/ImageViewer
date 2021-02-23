function updateColumns() {
  var elements = document.getElementsByTagName("a");
  for (var i = 0; i < elements.length; i++) {
    var url = elements[i].getAttribute("href");
    var columns = document.getElementById("number").value;
    console.log(url);
    console.log(url.replace(/\?columns=[0-9]*/, "?columns=" + columns));
    elements[i].setAttribute(
      "href",
      url.replace(/\?columns=[0-9]*/, "?columns=" + columns)
    );
  }
}
