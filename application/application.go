package application

import (
	"image/color"

	"github.com/adm87/finch-application/config"
	"github.com/adm87/finch-core/ecs"
	"github.com/adm87/finch-resources/storage"
	"github.com/hajimehoshi/ebiten/v2"
)

type ApplicationConfig struct {
	*config.Metadata
	*config.Resources
	*config.Window
}

type Application struct {
	config *ApplicationConfig

	startupFunc  func(app *Application) error // Called before the ebiten window is opened
	shutdownFunc func(app *Application) error

	cache *storage.ResourceCache
	world *ecs.World

	shouldExit bool
	clearColor color.RGBA
}

func NewApplication() *Application {
	return NewApplicationWithConfig(&ApplicationConfig{
		Metadata: &config.Metadata{
			Name:    "Finch",
			Version: "1.0.0",
			Root:    "/path/to/finch",
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
		config:     config,
		cache:      storage.NewResourceCache(),
		world:      ecs.NewWorld(),
		shouldExit: false,
		clearColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

func (app *Application) WithStartup(fn func(app *Application) error) *Application {
	app.startupFunc = fn
	return app
}

func (app *Application) WithShutdown(fn func(app *Application) error) *Application {
	app.shutdownFunc = fn
	return app
}

func (app *Application) Cache() *storage.ResourceCache {
	return app.cache
}

func (app *Application) World() *ecs.World {
	return app.world
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

// === Ebiten Game Interface ===

func (app *Application) Layout(outsideWidth, outsideHeight int) (int, int) {
	if window := app.Config().Window; window != nil {
		window.ScreenWidth = float32(outsideWidth) * window.RenderScale
		window.ScreenHeight = float32(outsideHeight) * window.RenderScale
		return int(window.ScreenWidth), int(window.ScreenHeight)
	}
	return outsideWidth, outsideHeight
}

func (app *Application) Draw(screen *ebiten.Image) {
	screen.Fill(app.clearColor)
}

func (app *Application) Update() error {
	if app.shouldExit {
		return ebiten.Termination
	}
	return nil
}
