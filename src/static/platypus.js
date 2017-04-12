window.onload = checkNight;
// Function that refreshes each
// panels stats indivudually
var table = document.getElementById("table");
var i = 0;
var ws = new WebSocket("ws://"+document.location.host+"/fetch");


ws.onmessage = function(evt) {
    //console.log(evt.data, panel)
    setRow(evt.data);
};


function setRow(text) {
    var res = JSON.parse(text);
    var panel = res["id"]
    console.log("---------- \n" + panel)
    var row = document.getElementById(panel);
    console.log(row.cells + "\n----------")
    //console.log(row);
    //console.log(panel)
    //console.log(res['online']);
    if (res['online'] == false || res['online'] == null) {
        row.cells[0].innerHTML = row.cells[0].innerHTML + " <strong>OFFLINE</strong>";
        row.cells[1].innerHTML = "0%";
        row.cells[2].innerHTML = "0%";
        row.cells[3].innerHTML = "0%";
    } else {    
        row.cells[1].innerHTML = res['cpu'] + "%";
        row.cells[2].innerHTML = res['memory'] + "%";
        row.cells[3].innerHTML = res['disk'] + "%";
    }
}

function toggleNight() {
    if (localStorage.nightmode != "true") {
        document.body.classList.add("nightmode");
        localStorage.setItem("nightmode", true);
    } else {
        document.body.classList.remove("nightmode");
        localStorage.setItem("nightmode", false);
    }
}

function checkNight() {
    if (localStorage.nightmode == "true") {
        console.log("Night time");
        document.getElementById("night").checked = true;
        document.body.classList.add("nightmode");
    } else {
        document.body.classList.remove("nightmode");
    }
}
