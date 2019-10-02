package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/structs"
)

func GetRoomConfig(ctx context.Context, client *http.Client, roomID string) (structs.Room, error) {
	var config structs.Room

	split := strings.Split(roomID, "-")
	if len(split) != 2 {
		return config, fmt.Errorf("invalid roomID %q, must be in the format BLDG-ROOM", roomID)
	}

	url := fmt.Sprintf("http://localhost:8000/buildings/%s/rooms/%s/configuration", split[0], split[1])
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return config, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return config, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(buf, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func GetRoomState(ctx context.Context, client *http.Client, roomID string) (structs.PublicRoom, error) {
	var state structs.PublicRoom

	split := strings.Split(roomID, "-")
	if len(split) != 2 {
		return state, fmt.Errorf("invalid roomID %q, must be in the format BLDG-ROOM", roomID)
	}

	url := fmt.Sprintf("http://localhost:8000/buildings/%s/rooms/%s", split[0], split[1])
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return state, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return state, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return state, err
	}

	err = json.Unmarshal(buf, &state)
	if err != nil {
		return state, err
	}

	return state, nil
}

func GetUIConfig(ctx context.Context, client *http.Client, roomID string) (structs.UIConfig, error) {
	return db.GetDB().GetUIConfig(roomID)
}
