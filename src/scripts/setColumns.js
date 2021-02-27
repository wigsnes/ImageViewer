function updateColumns() {
  var columns = document.getElementById("number").value;
  if (
    window.location.href.includes("?columns=" + columns) ||
    (!window.location.href.includes("?columns=") && columns === "3")
  ) {
    return;
  }

  var elements = document.getElementsByTagName("a");
  for (var i = 0; i < elements.length; i++) {
    var url = elements[i].getAttribute("href");
    elements[i].setAttribute(
      "href",
      url.replace(/\?columns=[0-9]*/, "?columns=" + columns)
    );
  }

  if (window.location.href.includes("?columns=")) {
    window.location.href = window.location.href.replace(
      /\?columns=[0-9]*/,
      "?columns=" + columns
    );
  } else {
    window.location.href = window.location.href + "?columns=" + columns;
  }
}
