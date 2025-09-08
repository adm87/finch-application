# finch-application
```go
package main

import (
	"image/color"

	finch "github.com/adm87/finch-application/application"
	"github.com/adm87/finch-application/config"
)

func main() {
	game := finch.NewApplicationWithConfig(
		&finch.ApplicationConfig{
			Metadata: &config.Metadata{
				Name:      "Finch Game",
				Root:      ".",
				TargetFps: 60,
			},
			Window: &config.Window{
				Title:           "Finch",
				Width:           800,
				Height:          600,
				ClearColor:      color.RGBA{R: 100, G: 149, B: 237, A: 255},
				ClearBackground: true,
			},
		},
	)

	cmd := finch.NewApplicationCommand("finch-game", game)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
```