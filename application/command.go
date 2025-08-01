package application

import (
	"path/filepath"

	"github.com/adm87/finch-resources/manifest"
	"github.com/spf13/cobra"
)

func NewApplicationCommand(use string, app *Application) *cobra.Command {
	cmd := &cobra.Command{
		Use: use,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(app.Config().Root)
			if err != nil {
				return err
			}
			app.Config().Root = root

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

				app.Cache().SetManifest(m)
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
		cmd.PersistentFlags().StringVar(&app.Config().Root, "root", ".", "Root directory for the application")
	}

	if window := app.Config().Window; window != nil {
		cmd.PersistentFlags().IntVar(&app.Config().Width, "window-width", 800, "Width of the application window")
		cmd.PersistentFlags().IntVar(&app.Config().Height, "window-height", 600, "Height of the application window")
		cmd.PersistentFlags().BoolVar(&app.Config().Fullscreen, "fullscreen", false, "Run the application in fullscreen mode")
	}

	if resources := app.Config().Resources; resources != nil {
		cmd.PersistentFlags().StringVar(&resources.Path, "resources-path", "data/", "Path to resources directory")
		cmd.PersistentFlags().StringVar(&resources.ManifestName, "manifest-name", "manifest.json", "Name of the resource manifest file")
	}

	return cmd
}
