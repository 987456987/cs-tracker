import { app, shell, BrowserWindow } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
const { spawn } = require('child_process')
import icon from '../../resources/icon.png?asset'

function createWindow() {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 900,
    height: 670,
    show: false,
    autoHideMenuBar: true,
    ...(process.platform === 'linux' ? { icon } : {}),
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false
    }
  })

  mainWindow.on('ready-to-show', () => {
    mainWindow.show()
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  // Set app user model id for windows
  electronApp.setAppUserModelId('com.electron')

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // see https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  createWindow()

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Spawn the Go backend as a subprocess
const path = require('path')
//DEV PATH
const combineExePath = path.join(app.getAppPath(), 'resources', 'combine.exe')
//BUILD PATH
// const combineExePath = path.join(
//   app.getAppPath(),
//   '..',
//   '..',
//   'resources',
//   'app.asar.unpacked',
//   'resources',
//   'combine.exe'
// )
const goBackendProcess = spawn(combineExePath, [
  /* any command-line arguments */
])

// Handle the output from the Go process
goBackendProcess.stdout.on('data', (data) => {
  console.log(`Go Backend Output: ${data}`)
  // You can send this data to your frontend if needed
})

goBackendProcess.stderr.on('data', (data) => {
  console.error(`Go Backend Error: ${data}`)
})

goBackendProcess.on('close', (code) => {
  console.log(`Go Backend exited with code ${code}`)
  // Delete .gz and .dem files before quitting the app
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    deleteFilesWithExtensions(['gz', 'dem'])
    console.log('HERE I AM')
    app.quit()
  }
})

// In this file you can include the rest of your app"s specific main process
// code. You can also put them in separate files and require them here.

const fs = require('fs')
//DELETE ALL FILES WITH EXTENSION
//MAY NEED TO CHANGE DIRECTORY ONCE BUILT
function deleteFilesWithExtensions(extensions) {
  const mainDirectory = app.getAppPath()
  fs.readdir(mainDirectory, (err, files) => {
    if (err) {
      console.error(`Error reading directory: ${err}`)
      return
    }

    files.forEach((file) => {
      const fileExtension = file.split('.').pop()
      if (extensions.includes(fileExtension)) {
        const filePath = path.join(mainDirectory, file)
        fs.unlink(filePath, (unlinkErr) => {
          if (unlinkErr) {
            console.error(`Error deleting file ${filePath}: ${unlinkErr}`)
          } else {
            console.log(`Deleted file: ${filePath}`)
          }
        })
      }
    })
  })
}
