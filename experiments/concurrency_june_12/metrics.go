// Remains the same

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const stateFile = "/app/latency_state.json"

type latencyState struct {
	MinStart      time.Time `json:"min_start"`
	MaxEnd        time.Time `json:"max_end"`
	TotalSeconds  float64   `json:"total_seconds"`
	Jobs          int       `json:"jobs"`
}

var state latencyState

func init() {
	// Try to load previous state (survives pod restarts if /app is persisted)
	data, err := os.ReadFile(stateFile)
	if err == nil {
		_ = json.Unmarshal(data, &state)
	}
}

func updateLatency(start, end time.Time) {
	delta := end.Sub(start).Seconds()
	log.Printf("Updated the latency")
	// first job ever
	if state.Jobs == 0 || start.Before(state.MinStart) {
		state.MinStart = start
	}
	if end.After(state.MaxEnd) {
		state.MaxEnd = end
	}
	// state.TotalSeconds += delta
	state.TotalSeconds = state.MaxEnd.Sub(state.MinStart).Seconds()
	state.Jobs++

	save()
}

func save() {
	file, err := os.CreateTemp("/app", "lat_tmp_*")
	if err != nil {
		log.Printf("latency save tmp err: %v", err)
		return
	}
	enc := json.NewEncoder(file)
	if err := enc.Encode(state); err != nil {
		log.Printf("latency encode err: %v", err)
	}
	file.Close()
	_ = os.Rename(file.Name(), stateFile) // atomic replace
}
