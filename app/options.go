package app

import (
	"io"
	"log/slog"
)

// ContainerLogsOpts struct to configure the container log
type ContainerLogsOpts struct {
	W io.Writer
	L slog.Leveler
}

// ContainerOpts application container options
type ContainerOpts struct {
	Log ContainerLogsOpts
}
