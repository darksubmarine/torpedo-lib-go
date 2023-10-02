package app

import (
	"io"
	"log/slog"
)

type ContainerLogsOpts struct {
	W io.Writer
	L slog.Leveler
}
type ContainerOpts struct {
	Log ContainerLogsOpts
}
