package router

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseEventOptions(t *testing.T) {
	t.Run("returns ids when repeated id query params are present", func(t *testing.T) {
		values := url.Values{}
		values.Add("id", "xyz1")
		values.Add("id", "xyz2")
		values.Add("startAfterId", "ignored")
		values.Add("limit", "50")

		got, err := ParseEventOptions(values)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}

		if !reflect.DeepEqual(got.IDs, []string{"xyz1", "xyz2"}) {
			t.Fatalf("expected ids [xyz1 xyz2] but got %v", got.IDs)
		}
		if got.Limit != 0 {
			t.Fatalf("expected limit to be ignored with ids but got %d", got.Limit)
		}
	})

	t.Run("uses pagination defaults when ids are missing", func(t *testing.T) {
		values := url.Values{}
		values.Add("startAfterId", "cursor-1")

		got, err := ParseEventOptions(values)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}

		if got.Limit != defaultEventLimit {
			t.Fatalf("expected default limit %d but got %d", defaultEventLimit, got.Limit)
		}
		if got.StartAfterID != "cursor-1" {
			t.Fatalf("expected startAfterID cursor-1 but got %q", got.StartAfterID)
		}
	})

	t.Run("validates limit", func(t *testing.T) {
		values := url.Values{}
		values.Add("limit", "0")

		_, err := ParseEventOptions(values)
		if err == nil {
			t.Fatal("expected error for invalid limit")
		}
	})
}

func TestParseSheetOptions(t *testing.T) {
	t.Run("returns repeated non-empty group keys", func(t *testing.T) {
		values := url.Values{}
		values.Add("group_key", "gos__carmen-habanera")
		values.Add("group_key", "")
		values.Add("group_key", "gos__suita")

		got := ParseSheetOptions(values)
		if !reflect.DeepEqual(got.GroupKeys, []string{"gos__carmen-habanera", "gos__suita"}) {
			t.Fatalf("expected group keys [gos__carmen-habanera gos__suita] but got %v", got.GroupKeys)
		}
	})
}
