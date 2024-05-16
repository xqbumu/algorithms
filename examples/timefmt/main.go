package main

import (
	"log/slog"
	"time"
)

func main() {
	now := time.Now()
	slog.Info(now.Format(time.RFC1123))
}
