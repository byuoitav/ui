package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/byuoitav/common/structs"
)

func GetRoomConfig(ctx context.Context, client *http.Client, roomID string) (structs.Room, error) {
	var config structs.Room

	split := strings.Split(roomID, "-")
	if len(split) != 2 {
		return config, fmt.Errorf("invalid roomID %q, must be in the format BLDG-ROOM", roomID)
	}

	// TODO use the one in aws
	url := fmt.Sprintf("http://itb-1006-cp1.byu.edu:8000/buildings/%s/rooms/%s/configuration", split[0], split[1])
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

	// TODO use the one in aws
	url := fmt.Sprintf("http://itb-1006-cp1.byu.edu:8000/buildings/%s/rooms/%s", split[0], split[1])
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

func GetUIConfig(ctx context.Context, client *http.Client, roomID string) (UIConfig, error) {
	var config UIConfig

	endpoint := fmt.Sprintf("ui-configuration/%s", roomID)
	url := fmt.Sprintf("%s/%s", os.Getenv("DB_ADDRESS"), endpoint)
	url = strings.TrimSpace(url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return config, err
	}

	uname := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	if len(uname) > 0 && len(pass) > 0 {
		req.SetBasicAuth(uname, pass)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return config, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse response: %s. response: %s", err, b)
	}

	if len(config.ID) == 0 {
		return config, fmt.Errorf("unable to get %s: %s", endpoint, b)
	}

	return config, nil
}
