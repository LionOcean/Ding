package main

import (
	"context"
	"ding/transfer"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// ---- app lifeCycle ----

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	runtime.WindowShow(a.ctx)
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// ---- app native go methods bindings to js ----

// StartP2PServer start a server for received perr
func (a *App) StartP2PServer() error {
	return transfer.StartP2PServer()
}

// CloseP2PServer close the runing P2PServer
func (a *App) CloseP2PServer() error {
	return transfer.CloseP2PServer()
}

// LocalIPAddr return current peer local ipv4 and port, if in received peer, you should omit port filed.
func (a *App) LocalIPAddr() ([]string, error) {
	return transfer.LocalIPAddr()
}

// UploadFiles show a system dialog to choose files and upload to server
func (a *App) UploadFiles(dialogOptions runtime.OpenDialogOptions) ([]transfer.TransferFile, error) {
	emptyfiles := make([]transfer.TransferFile, 0)
	paths, err := runtime.OpenMultipleFilesDialog(a.ctx, dialogOptions)
	if err != nil {
		return emptyfiles, err
	}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return emptyfiles, err
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			return emptyfiles, err
		}
		transfer.AppendTransferFile(transfer.TransferFile{
			Path: path,
			Name: fileInfo.Name(),
			Size: int(fileInfo.Size()),
		})
	}
	return transfer.LogTransferFiles(), nil
}

// RemoveTransferFiles remove files from files_list.
func (a *App) RemoveFiles(files []transfer.TransferFile) {
	transfer.RemoveTransferFiles(files...)
}

// SaveFileDialog show a system dialog to choose a saving file path
func (a *App) SaveFileDialog(dialogOptions runtime.SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(a.ctx, dialogOptions)
}

// DownloadFile make a GET request to remotePath, next write response.body to local file with buffer pieces.
func (a *App) DownloadFile(remotePath, localPath string) error {
	return transfer.DownloadFile(remotePath, localPath)
}
