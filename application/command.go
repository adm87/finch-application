package application

import (
	"path"
	"path/filepath"

	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/fsys"
	"github.com/adm87/finch-resources/manifest"
	"github.com/adm87/finch-resources/resources"
	"github.com/spf13/cobra"
)

func NewApplicationCommand(use string, app *Application) *cobra.Command {
	cmd := &cobra.Command{
		Use: use,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if metadata := app.Config().Metadata; metadata == nil {
				return errors.NewInvalidArgumentError("application metadata must be set")
			}

			root, err := filepath.Abs(app.Config().Root)
			if err != nil {
				return err
			}
			app.Config().Root = filepath.ToSlash(root)

			if window := app.Config().Window; window != nil {
				if window.Width <= 0 || window.Height <= 0 {
					return errors.NewInvalidArgumentError("window width and height must be non-zero positive integers")
				}
				if window.RenderScale <= 0 {
					window.RenderScale = 1.0
				}
				window.ScreenWidth = float32(window.Width) * window.RenderScale
				window.ScreenHeight = float32(window.Height) * window.RenderScale
			}

			if res := app.Config().Resources; res != nil {
				res.Path = path.Join(root, res.Path)

				if res.ManifestName == "" {
					res.ManifestName = "manifest.json"
				}

				manifestPath := path.Join(res.Path, res.ManifestName)

				if fsys.Exists(manifestPath) {
					m, err := manifest.LoadManifest(manifestPath)

					if err != nil {
						return err
					}

					resources.SetManifest(m)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Open()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	if metadata := app.Config().Metadata; metadata != nil {
		cmd.PersistentFlags().StringVar(&metadata.Root, "root", metadata.Root, "Root directory for the application")
	}

	if window := app.Config().Window; window != nil {
		cmd.PersistentFlags().IntVar(&window.Width, "window-width", window.Width, "Width of the application window")
		cmd.PersistentFlags().IntVar(&window.Height, "window-height", window.Height, "Height of the application window")
		cmd.PersistentFlags().BoolVar(&window.Fullscreen, "fullscreen", window.Fullscreen, "Run the application in fullscreen mode")
	}

	if res := app.Config().Resources; res != nil {
		cmd.PersistentFlags().StringVar(&res.Path, "resources-path", res.Path, "Path to resources directory")
		cmd.PersistentFlags().StringVar(&res.ManifestName, "manifest-name", res.ManifestName, "Name of the resource manifest file")
	}

	return cmd
}
