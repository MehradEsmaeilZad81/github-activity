package activity

import "github-activity/internal/ghapi"

func FallbackLine(e ghapi.EventRaw) string {
	if e.Repo.Name != "" && e.Type != "" {
		return e.Type + " on " + e.Repo.Name
	}
	if e.Type != "" {
		return e.Type
	}
	return ""
}
