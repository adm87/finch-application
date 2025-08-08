package time

import (
	stdtime "time"

	"github.com/adm87/finch-core/errors"
)

type FPS struct {
	targetMs          float64
	deltaTime         float64
	elapsedTime       float64
	startTime         float64
	currentTime       float64
	deltaSeconds      float64
	fixedDeltaSeconds float64
	lastFpsUpdate     float64
	fps               float64
	frameCount        int
}

func NewFPS(target int) *FPS {
	if target <= 0 {
		panic(errors.NewInvalidArgumentError("target FPS must be greater than 0"))
	}
	targetMs := 1000.0 / float64(target)
	return &FPS{
		targetMs:          targetMs,
		fixedDeltaSeconds: targetMs / 1000.0,
	}
}

func (f *FPS) Start() {
	if f.startTime != 0 {
		return
	}

	now := float64(stdtime.Now().UnixMilli())

	f.startTime = now
	f.currentTime = now
	f.lastFpsUpdate = now
	f.fps = 0
	f.elapsedTime = 0
	f.deltaTime = 0
	f.deltaSeconds = 0
	f.frameCount = 0
}

// Update calculates the FPS based on the elapsed time.
//
// It returns the number of fixed frames that occurred since the last update.
func (f *FPS) Update() (frames int) {
	if f.startTime == 0 {
		f.Start()
	}

	now := float64(stdtime.Now().UnixMilli())
	previousTime := f.currentTime

	f.currentTime = now
	f.deltaTime = f.currentTime - previousTime

	f.elapsedTime += f.deltaTime
	f.deltaSeconds = float64(f.deltaTime) / 1000.0

	frames = int(f.elapsedTime / f.targetMs)
	if frames > 0 {
		f.elapsedTime -= float64(frames) * f.targetMs
	}

	f.frameCount += frames

	fpsUpdateElapsed := now - f.lastFpsUpdate
	if fpsUpdateElapsed >= 1000 {
		f.fps = float64(f.frameCount) / (float64(fpsUpdateElapsed) / 1000.0)
		f.frameCount = 0
		f.lastFpsUpdate = now
	}

	return frames
}

func (f *FPS) DeltaSeconds() float64 {
	return f.deltaSeconds
}

func (f *FPS) FixedDeltaSeconds() float64 {
	return f.fixedDeltaSeconds
}

// Interpolation calculates the interpolation factor for the current frame.
func (f *FPS) Interpolation() float64 {
	now := float64(stdtime.Now().UnixMilli())

	elapsedSinceUpdate := now - f.currentTime
	interpolation := (f.elapsedTime + elapsedSinceUpdate) / f.targetMs

	if interpolation > 1.0 {
		interpolation = 1.0
	} else if interpolation < 0.0 {
		interpolation = 0.0
	}

	return interpolation
}

func (f *FPS) FPS() float64 {
	return f.fps
}
