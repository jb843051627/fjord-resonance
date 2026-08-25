package store

import (
	"fmt"
	"strings"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func AuditMessage(event model.AuditEvent) string {
	parts := []string{event.Action, string(event.Entity), string(event.EntityID)}
	if strings.TrimSpace(event.Actor) != "" {
		parts = append(parts, "by", event.Actor)
	}
	if strings.TrimSpace(event.Details) != "" {
		parts = append(parts, event.Details)
	}
	return fmt.Sprintf("%s", strings.Join(parts, " "))
}
