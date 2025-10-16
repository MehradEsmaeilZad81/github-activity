package activity

import (
	"encoding/json"
	"fmt"
	"github-activity/internal/ghapi"
)

type PushPayload struct {
	Size int `json:"size"`
}

type IssuesPayload struct {
	Action string `json:"action"`
	Issue  struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	} `json:"issue"`
}

type WatchPayload struct {
	Action string `json:"action"`
}

type PullRequestPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	} `json:"pull_request"`
}

type CreatePayload struct {
	RefType string `json:"ref_type"`
	Ref     string `json:"ref"`
}

type ForkPayload struct {
	Forkee struct {
		FullName string `json:"full_name"`
	} `json:"forkee"`
}

type IssueCommentPayload struct {
	Action string `json:"action"`
	Issue  struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	} `json:"issue"`
}

type ReleasePayload struct {
	Action  string `json:"action"`
	Release struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
	} `json:"release"`
}

func FormatEvent(ev ghapi.EventRaw) string {
	switch ev.Type {
	case "PushEvent":
		var p PushPayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			if p.Size == 1 {
				return fmt.Sprintf("Pushed %d commit to %s", p.Size, ev.Repo.Name)
			}
			return fmt.Sprintf("Pushed %d commits to %s", p.Size, ev.Repo.Name)
		}
	case "IssuesEvent":
		var p IssuesPayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			return fmt.Sprintf("%s issue #%d in %s — %s",
				capitalize(p.Action), p.Issue.Number, ev.Repo.Name, p.Issue.Title)
		}
	case "WatchEvent": // star
		var p WatchPayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			return fmt.Sprintf("%s %s", capitalize(p.Action), ev.Repo.Name)
		}
	case "PullRequestEvent":
		var p PullRequestPayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			return fmt.Sprintf("%s pull request #%d in %s — %s",
				capitalize(p.Action), p.PullRequest.Number, ev.Repo.Name, p.PullRequest.Title)
		}
	case "CreateEvent":
		var p CreatePayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			if p.RefType == "repository" {
				return fmt.Sprintf("Created repository %s", ev.Repo.Name)
			}
			if p.Ref != "" {
				return fmt.Sprintf("Created %s %s in %s", p.RefType, p.Ref, ev.Repo.Name)
			}
			return fmt.Sprintf("Created %s in %s", p.RefType, ev.Repo.Name)
		}
	case "ForkEvent":
		var p ForkPayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			if p.Forkee.FullName != "" {
				return fmt.Sprintf("Forked %s to %s", ev.Repo.Name, p.Forkee.FullName)
			}
			return fmt.Sprintf("Forked %s", ev.Repo.Name)
		}
	case "IssueCommentEvent":
		var p IssueCommentPayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			return fmt.Sprintf("%s a comment on issue #%d in %s — %s",
				capitalize(p.Action), p.Issue.Number, ev.Repo.Name, p.Issue.Title)
		}
	case "ReleaseEvent":
		var p ReleasePayload
		if json.Unmarshal(ev.Payload, &p) == nil {
			if p.Release.Name != "" {
				return fmt.Sprintf("%s release %s (%s) in %s",
					capitalize(p.Action), p.Release.TagName, p.Release.Name, ev.Repo.Name)
			}
			return fmt.Sprintf("%s release %s in %s",
				capitalize(p.Action), p.Release.TagName, ev.Repo.Name)
		}
	default:
		return FallbackLine(ev)
	}
	return ""
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] = r[0] - ('a' - 'A')
	}
	return string(r)
}
