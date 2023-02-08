package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

// OpenFileDialog show a system dialog to choose files
func (a *App) OpenFileDialog(dialogOptions runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenFileDialog(a.ctx, dialogOptions)
}

// SaveFileDialog show a system dialog to choose a saving file path
func (a *App) SaveFileDialog(dialogOptions runtime.SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(a.ctx, dialogOptions)
}

// WriteFile write string data to path
func (a *App) WriteFile(path string, data string) error {
	// 创建文件
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入buf到文件
	_, writeErr := file.Write([]byte(data))
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// Extname return path ext name
func (a *App) Extname(path string) string {
	return filepath.Ext(path)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
