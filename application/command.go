package application

import (
	"path/filepath"

	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-resources/manifest"
	"github.com/adm87/finch-resources/storage"
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
			app.Config().Root = root

			if window := app.Config().Window; window != nil {
				if window.Width <= 0 || window.Height <= 0 {
					return errors.NewInvalidArgumentError("window width and height must be non-zero positive integers")
				}
				if window.RenderScale <= 0 {
					return errors.NewInvalidArgumentError("window render scale must be a non-zero positive integer")
				}

				window.ScreenWidth = float32(window.Width) * window.RenderScale
				window.ScreenHeight = float32(window.Height) * window.RenderScale
			}

			if resources := app.Config().Resources; resources != nil {
				resources.Path = filepath.Join(root, resources.Path)

				if resources.ManifestName == "" {
					resources.ManifestName = "manifest.json"
				}

				manifestPath := filepath.Join(resources.Path, resources.ManifestName)
				m, err := manifest.LoadManifest(manifestPath)

				if err != nil {
					return err
				}

				storage.SetManifest(m)
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

	if resources := app.Config().Resources; resources != nil {
		cmd.PersistentFlags().StringVar(&resources.Path, "resources-path", resources.Path, "Path to resources directory")
		cmd.PersistentFlags().StringVar(&resources.ManifestName, "manifest-name", resources.ManifestName, "Name of the resource manifest file")
	}

	return cmd
}
