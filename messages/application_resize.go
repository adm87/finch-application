package messages

import (
	"github.com/adm87/finch-core/geometry"
	"github.com/adm87/finch-core/messaging"
)

type ApplicationResizeMessage struct {
	To   geometry.Point
	From geometry.Point
}

var ApplicationResize = messaging.NewMessageBus[ApplicationResizeMessage]()
