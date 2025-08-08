package config

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Window holds configuration for the application window.
type Window struct {
	Title           string
	Width           int
	Height          int
	ScreenWidth     float32
	ScreenHeight    float32
	ResizeMode      ebiten.WindowResizingModeType
	RenderScale     float32
	Fullscreen      bool
	ClearBackground bool
	ClearColor      color.RGBA
}
