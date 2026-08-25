package api

import (
	"net/url"
	"strconv"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func BatchFilterFromQuery(values url.Values) store.BatchFilter {
	filter := store.BatchFilter{BuoyID: modelIDQuery(values.Get("buoy_id")), Status: statusQuery(values.Get("status")), Limit: intQuery(values.Get("limit"), 50), Offset: intQuery(values.Get("offset"), 0)}
	if value := values.Get("from"); value != "" {
		filter.From, _ = time.Parse(time.RFC3339, value)
	}
	if value := values.Get("to"); value != "" {
		filter.To, _ = time.Parse(time.RFC3339, value)
	}
	return filter
}

func intQuery(value string, fallback int) int {
	number, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return number
}

func modelIDQuery(value string) model.ID { return model.ID(value) }

func statusQuery(value string) model.BatchStatus { return model.BatchStatus(value) }
