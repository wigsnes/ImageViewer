refresh_handler = function(e) {
    var elements = document.querySelectorAll("*[data-src]");
    for (var i = 0; i < elements.length; i++) {
        var boundingClientRect = elements[i].getBoundingClientRect();
        if (elements[i].hasAttribute("data-src") && boundingClientRect.top < window.innerHeight) {
            if (elements[i].tagName === "IMG") {
                elements[i].setAttribute("src", elements[i].getAttribute("data-src"));
                elements[i].removeAttribute("data-src");
                setTimeout(function(elem) {
                    elem.setAttribute("style", "height: auto;");
                }, 100, elements[i]);
            } else {
                elements[i].innerHTML = "<source src='" + elements[i].getAttribute("data-src") + "'type='" + elements[i].getAttribute("data-type") + "' />"
                elements[i].removeAttribute("data-src");
                elements[i].removeAttribute("data-type");
            }
        }
    }
};
wait_refresh_handler = function(e) {
    setTimeout(
        function() {
            refresh_handler();
        }, 1000
    );
};

window.addEventListener('scroll', refresh_handler);
window.addEventListener('load', wait_refresh_handler);
window.onload = refresh_handler;