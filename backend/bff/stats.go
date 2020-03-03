package bff

import "time"

type ClientStats struct {
	CreatedAt *time.Time `json:"CreatedAt,omitempty"`
	Routines  uint

	WebSocket struct {
		MessagesRecieved uint
		MessagesSent     uint
		ErrorsSent       uint
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
	ClientCount  uint
	OldestClient *time.Time `json:"OldestClient,omitempty"`

	ClientStats
}

func (c *Client) Stats() ClientStats {
	stats := c.stats

	return stats
}

func AggregateStats(stats []ClientStats) AggregateClientStats {
	agg := AggregateClientStats{}
	agg.AvControlApi.ResponseCodes = make(map[int]uint)

	for i := range stats {
		agg.ClientCount++

		// figure out the oldest client
		switch {
		case stats[i].CreatedAt == nil:
		case agg.OldestClient == nil:
			agg.OldestClient = stats[i].CreatedAt
		case stats[i].CreatedAt.Before(*agg.OldestClient):
			agg.OldestClient = stats[i].CreatedAt
		}

		agg.Routines += stats[i].Routines

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

func (s *ClientStats) decRoutines() {
	s.Routines--
}
