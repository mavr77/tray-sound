package main

import (
	"encoding/json"
	"testing"
)

const sampleJSON = `{"slug":"kyiv","type":"state","name_uk":"м. Київ","name_en":"Kyiv","active":true,"level":"yellow","started_at":"2026-11-02T15:09:07.000000Z","updated_at":"2026-11-02T15:42:50.705565Z","districts_active":[],"active_anywhere":true}`

func TestParseSample(t *testing.T) {
	var st alertStatus
	if err := json.Unmarshal([]byte(sampleJSON), &st); err != nil {
		t.Fatal(err)
	}
	if !st.Active || st.level() != "yellow" || st.color() != colYellow {
		t.Fatalf("unexpected: %+v", st)
	}
	if st.StartedAt.IsZero() {
		t.Fatal("started_at not parsed")
	}

	st.Active = false
	if st.color() != colGreen || st.key() != "off" {
		t.Fatal("inactive should be green")
	}

	st.Active, st.Level = true, "something-new"
	if st.color() != colorForUnknownLevel {
		t.Fatal("unknown level should use fallback color")
	}
}
