package application

import (
	"image/color"

	"github.com/adm87/finch-application/config"
	"github.com/adm87/finch-application/time"
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
	fps   *time.FPS

	shouldExit bool
	clearColor color.RGBA
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
		config:     config,
		cache:      storage.NewResourceCache(),
		world:      ecs.NewWorld(),
		fps:        time.NewFPS(config.TargetFps),
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
		return int(window.ScreenWidth), int(window.ScreenHeight)
	}

	return outsideWidth, outsideHeight
}

func (app *Application) Draw(screen *ebiten.Image) {
	if window := app.Config().Window; window != nil && window.ClearBackground {
		screen.Fill(window.ClearColor)
	}
	app.world.Render(screen, app.fps.Interpolation())
}

func (app *Application) Update() error {
	if app.shouldExit {
		return app.internal_shutdown()
	}

	fixedFrames := app.fps.Update()

	if err := app.world.EarlyUpdate(app.fps.DeltaSeconds()); err != nil {
		return err
	}

	if fixedFrames > 0 {
		const maxFixedUpdates = 5
		if fixedFrames > maxFixedUpdates {
			fixedFrames = maxFixedUpdates
		}

		fixedDeltaSeconds := app.fps.FixedDeltaSeconds()
		for i := 0; i < fixedFrames; i++ {
			if err := app.world.FixedUpdate(fixedDeltaSeconds); err != nil {
				return err
			}
		}
	}

	if err := app.world.LateUpdate(app.fps.DeltaSeconds()); err != nil {
		return err
	}

	return nil
}
