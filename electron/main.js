const electron = require("electron");
const { dialog } = require("electron");
var XMLHttpRequest = require("xmlhttprequest").XMLHttpRequest;
// Module to control application life.
const app = electron.app;
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow;

const exec = require("child_process").exec;
const path = require("path");
const url = require("url");

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow;

function execute(command, callback) {
  exec(command, (error, stdout, stderr) => {
    callback(stdout);
  });
}

let options = {
  title: "Open directory",
  properties: ["openDirectory"],
};

async function createWindow() {
  // Create the browser window.
  mainWindow = new BrowserWindow({
    webPreferences: {
      webSecurity: false,
    },
  });

  mainWindow.webContents.session.clearCache(function () {
    //some callback.
  });

  // let filePath = await dialog.showOpenDialog(mainWindow, options);

  // execute('go run ../src/main.go -path="' + filePath[0] + "\\", (output) => {
  //   console.log(output);
  // });

  mainWindow.loadURL("http://localhost:8080");

  // Emitted when the window is closed.
  mainWindow.on("closed", function () {
    const Http = new XMLHttpRequest();
    const url = "http://localhost:8080/exit";
    Http.open("GET", url);
    Http.send();
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null;
  });
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", createWindow);

// Quit when all windows are closed.
app.on("window-all-closed", function () {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== "darwin") {
    app.quit();
  }
});

app.on("activate", function () {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    createWindow();
  }
});
