package cli

import (
	"context"
	"fmt"
	"github-activity/internal/activity"
	"github-activity/internal/ghapi"
	"io"
	"os"
	"time"
)

type Runner struct {
	API ghapi.Client
}

func (r Runner) Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) < 1 {
		_, _ = fmt.Fprintln(stderr, "Usage: github-activity <github-username>")
		return 1
	}
	username := args[0]

	if r.API == nil {
		_, _ = fmt.Fprintln(stderr, "internal error: API client is nil")
		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := r.API.UserEvents(ctx, username)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	if len(events) == 0 {
		_, _ = fmt.Fprintln(stdout, "No recent public activity found.")
		return 0
	}
	for _, ev := range events {
		line := activity.FormatEvent(ev)
		if line == "" {
			line = activity.FallbackLine(ev)
		}
		if line != "" {
			_, _ = fmt.Fprintln(stdout, "-", line)
		}
	}
	return 0
}

func Main(api ghapi.Client) int {
	r := Runner{API: api}
	return r.Run(os.Args[1:], os.Stdout, os.Stderr)
}
