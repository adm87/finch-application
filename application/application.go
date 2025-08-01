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

	startupFunc  func(app *Application) error
	shutdownFunc func(app *Application) error

	cache *storage.ResourceCache
	world *ecs.World

	shouldExit bool
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
