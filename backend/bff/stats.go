package bff

import "time"

type ClientStats struct {
	CreatedAt *time.Time `json:"CreatedAt,omitempty"`

	WebSocket struct {
		MessagesRecieved uint
		MessagesSent     uint
		ErrorsSent       uint
		TimeOpen         time.Duration `json:"TimeOpen,omitempty"`
	}

	AvControlApi struct {
		Requests      uint
		ResponseCodes map[int]uint
	}

	Lazarette struct {
		UpdatesRecieved uint
		UpdatesSent     uint
	}

	Events struct {
		Recieved uint
		Sent     uint
	}
}

type AggregateClientStats struct {
	ClientCount uint
	ClientStats
}

func (c *Client) Stats() ClientStats {
	stats := c.stats

	if stats.CreatedAt != nil {
		stats.WebSocket.TimeOpen = time.Since(*stats.CreatedAt)
	}

	return stats
}

func AggregateStats(stats []ClientStats) AggregateClientStats {
	agg := AggregateClientStats{}
	agg.AvControlApi.ResponseCodes = make(map[int]uint)

	for i := range stats {
		agg.ClientCount++

		agg.WebSocket.MessagesRecieved += stats[i].WebSocket.MessagesRecieved
		agg.WebSocket.MessagesSent += stats[i].WebSocket.MessagesSent
		agg.WebSocket.ErrorsSent += stats[i].WebSocket.ErrorsSent

		agg.AvControlApi.Requests += stats[i].AvControlApi.Requests
		for code, count := range stats[i].AvControlApi.ResponseCodes {
			agg.AvControlApi.ResponseCodes[code] += count
		}

		agg.Lazarette.UpdatesRecieved += stats[i].Lazarette.UpdatesRecieved
		agg.Lazarette.UpdatesSent += stats[i].Lazarette.UpdatesSent

		agg.Events.Recieved += stats[i].Events.Recieved
		agg.Events.Sent += stats[i].Events.Sent
	}

	return agg
}
