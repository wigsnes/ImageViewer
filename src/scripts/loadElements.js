refresh_handler = function (e) {
  var elements = document.querySelectorAll("*[data-src]");
  for (var i = 0; i < elements.length; i++) {
    var boundingClientRect = elements[i].getBoundingClientRect();
    if (
      elements[i].hasAttribute("data-src") &&
      boundingClientRect.top < window.innerHeight
    ) {
      elements[i].setAttribute("src", elements[i].getAttribute("data-src"));
      elements[i].removeAttribute("data-src");
      setTimeout(
        function (elem) {
          elem.setAttribute("style", "height: auto;");
        },
        100,
        elements[i]
      );
    }
  }
};

wait_refresh_handler = function (e) {
  setTimeout(function () {
    refresh_handler();
  }, 1000);
};

let waiting = false;

window.onload = function () {
  refresh_handler();

  document.getElementById("page").onscroll = function () {
    if (waiting) {
      return;
    }
    waiting = true;

    refresh_handler();

    setTimeout(function () {
      waiting = false;
    }, 100);
  };

  document
    .getElementById("page")
    .addEventListener("load", wait_refresh_handler);
};
