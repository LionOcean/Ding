package main

import (
	"context"
	"ding/transfer"

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
	// need wait for SPA web javascript rendering
	// runtime.WindowShow(a.ctx)
}
func (a *App) LocalIPv4s() ([]string, error) {
	return transfer.LocalIPv4s()
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

// ---- wails runtime native go methods bindings to js ----

// SaveFileDialog show a system dialog to choose a saving file path
func (a *App) SaveFileDialog(dialogOptions runtime.SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(a.ctx, dialogOptions)
}

// OpenDirDialog show a system dialog to choose a saving directory
func (a *App) OpenDirDialog(dialogOptions runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, dialogOptions)
}
