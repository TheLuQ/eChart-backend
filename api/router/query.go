package router

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	defaultEventLimit = 20
	minEventLimit     = 1
	maxEventLimit     = 100
)

type EventOptions struct {
	IDs          []string
	Limit        int
	StartAfterID string
}

type SheetOptions struct {
	GroupKeys []string
}

func ParseEventOptions(values url.Values) (EventOptions, error) {
	ids := nonEmpty(values["id"])
	if len(ids) > 0 {
		return EventOptions{IDs: ids}, nil
	}

	limit := defaultEventLimit
	limitQuery := values.Get("limit")
	if limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil || parsedLimit < minEventLimit || parsedLimit > maxEventLimit {
			return EventOptions{}, fmt.Errorf("limit must be an integer between %d and %d", minEventLimit, maxEventLimit)
		}
		limit = parsedLimit
	}

	return EventOptions{
		Limit:        limit,
		StartAfterID: values.Get("startAfterId"),
	}, nil
}

func ParseSheetOptions(values url.Values) SheetOptions {
	return SheetOptions{GroupKeys: nonEmpty(values["group_key"])}
}

func nonEmpty(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}
