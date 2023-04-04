package main

import (
	"github.com/gookit/slog"
	"github.com/robbailey3/chat-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error(err)
	}
}
