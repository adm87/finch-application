package application

import (
	"image/color"

	"github.com/adm87/finch-application/config"
	"github.com/adm87/finch-application/messages"
	"github.com/adm87/finch-application/time"
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-resources/storage"
	"github.com/hajimehoshi/ebiten/v2"
)

type StartupShutdownFunc func(app *Application) error
type DrawFunc func(app *Application, screen *ebiten.Image) error
type UpdateFunc func(app *Application, dt, fdt float64, frames int) error

type ApplicationConfig struct {
	*config.Metadata
	*config.Resources
	*config.Window
}

type Application struct {
	config *ApplicationConfig

	startupFunc  StartupShutdownFunc
	shutdownFunc StartupShutdownFunc
	drawFunc     DrawFunc
	updateFunc   UpdateFunc

	cache *storage.ResourceCache
	fps   *time.FPS

	shouldExit bool
	clearColor color.RGBA

	activeWidth  float32
	activeHeight float32
}

func NewApplication() *Application {
	return NewApplicationWithConfig(&ApplicationConfig{
		Metadata: &config.Metadata{
			Name:      "Finch",
			Version:   "1.0.0",
			Root:      "/path/to/finch",
			TargetFps: 30,
		},
		Resources: &config.Resources{
			Path:         "/path/to/resources",
			ManifestName: "manifest.json",
		},
		Window: &config.Window{
			Title:        "Finch Application",
			Width:        800,
			Height:       600,
			ScreenWidth:  800,
			ScreenHeight: 600,
			ResizeMode:   ebiten.WindowResizingModeEnabled, // Default resizing mode
			RenderScale:  1.0,
			Fullscreen:   false,
			ClearColor:   color.RGBA{R: 100, G: 149, B: 237, A: 255}, // Cornflower Blue
		},
	})
}

func NewApplicationWithConfig(config *ApplicationConfig) *Application {
	return &Application{
		config:       config,
		cache:        storage.NewResourceCache(),
		fps:          time.NewFPS(config.TargetFps, 5),
		shouldExit:   false,
		clearColor:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
		activeWidth:  0,
		activeHeight: 0,
	}
}

func (app *Application) WithStartup(fn StartupShutdownFunc) *Application {
	app.startupFunc = fn
	return app
}

func (app *Application) WithShutdown(fn StartupShutdownFunc) *Application {
	app.shutdownFunc = fn
	return app
}

func (app *Application) WithDraw(fn DrawFunc) *Application {
	app.drawFunc = fn
	return app
}

func (app *Application) WithUpdate(fn UpdateFunc) *Application {
	app.updateFunc = fn
	return app
}

func (app *Application) Cache() *storage.ResourceCache {
	return app.cache
}

func (app *Application) Config() *ApplicationConfig {
	return app.config
}

func (app *Application) Open() error {
	if window := app.Config().Window; window != nil {
		ebiten.SetWindowTitle(window.Title)
		ebiten.SetWindowSize(window.Width, window.Height)
		ebiten.SetWindowResizingMode(window.ResizeMode)
		ebiten.SetFullscreen(window.Fullscreen)

		app.clearColor = window.ClearColor
	}
	if app.startupFunc != nil {
		if err := app.startupFunc(app); err != nil {
			return err
		}
	}
	return ebiten.RunGame(app)
}

func (app *Application) Quit() {
	app.shouldExit = true
}

func (app *Application) internal_shutdown() error {
	if app.shutdownFunc != nil {
		if err := app.shutdownFunc(app); err != nil {
			return err
		}
	}
	return ebiten.Termination
}

// === Ebiten Game Interface ===

func (app *Application) Layout(outsideWidth, outsideHeight int) (int, int) {
	if window := app.Config().Window; window != nil {
		window.ScreenWidth = float32(outsideWidth) * window.RenderScale
		window.ScreenHeight = float32(outsideHeight) * window.RenderScale

		if app.activeWidth != window.ScreenWidth || app.activeHeight != window.ScreenHeight {
			fromWidth := app.activeWidth
			fromHeight := app.activeHeight

			app.activeWidth = window.ScreenWidth
			app.activeHeight = window.ScreenHeight

			messages.ApplicationResize.Send(messages.ApplicationResizeMessage{
				To:   geometry.Point{X: window.ScreenWidth, Y: window.ScreenHeight},
				From: geometry.Point{X: fromWidth, Y: fromHeight},
			})
		}

		return int(window.ScreenWidth), int(window.ScreenHeight)
	}
	return outsideWidth, outsideHeight
}

func (app *Application) Draw(screen *ebiten.Image) {
	if window := app.Config().Window; window != nil && window.ClearBackground {
		screen.Fill(window.ClearColor)
	}

	if app.drawFunc != nil {
		if err := app.drawFunc(app, screen); err != nil {
			panic(err)
		}
	}

}

func (app *Application) Update() error {
	if app.shouldExit {
		return app.internal_shutdown()
	}

	deltaSeconds, fixedDeltaSeconds, frames := app.fps.Update()

	if app.updateFunc != nil {
		if err := app.updateFunc(app, deltaSeconds, fixedDeltaSeconds, frames); err != nil {
			return err
		}
	}

	return nil
}
