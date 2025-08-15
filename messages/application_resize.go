package messages

import (
	"github.com/adm87/finch-core/events"
	"github.com/adm87/finch-core/geometry"
)

type ApplicationResizeMessage struct {
	To   geometry.Point
	From geometry.Point
}

var ApplicationResize = events.NewMessageBus[ApplicationResizeMessage]()
