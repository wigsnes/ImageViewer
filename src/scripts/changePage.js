function prevPage(page) {
    const Http = new XMLHttpRequest();
    const url="http://localhost:8080?page=" + page;
    Http.open("GET", url);
    Http.send();
}

function nextPage(page) {
    const Http = new XMLHttpRequest();
    const url="http://localhost:8080?page=" + page;
    Http.open("GET", url);
    Http.send();
}