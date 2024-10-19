package lifecycle

import "github.com/4strodev/rabbit_go_publisher/pkg/core/components"

type OnShutdown interface {
	components.Component
	OnShutdown() error
}
