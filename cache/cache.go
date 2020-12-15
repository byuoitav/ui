package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/ui"
	bolt "go.etcd.io/bbolt"
)

const (
	_configBucket              = "configs"
	_uiForDeviceBucket         = "uiForDevice"
	_controlGroupBucket        = "controlGroups"
	_roomAndControlGroupBucket = "roomAndControlGroup"
)

type dataService struct {
	dataService ui.DataService
	db          *bolt.DB
}

type roomAndControlgroup struct {
	Room         string
	ControlGroup string
}

func New(ds ui.DataService, path string) (ui.DataService, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to open cache: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(_uiForDeviceBucket))
		if err != nil {
			return err
		}

		_, err := tx.CreateBucketIfNotExists([]byte(_configBucket))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cache: %w", err)
	}

	return &dataService{
		dataService: ds,
		db:          db,
	}, nil
}

func (d *dataService) UIForDevice(ctx context.Context, room, id string) (string, error) {
	uiForDev, err := d.dataService.UIForDevice(ctx, room, id)
	if err != nil {
		uiForDev, cacheErr := d.uiForDeviceFromCache(ctx, room, id)
		if cacheErr != nil {
			log.L.Warnf("unable to get ui for device %s %s from cache: %s", room, id, cacheErr)
			return uiForDev, err
		}

		return uiForDev, nil
	}

	if err := d.cacheUIForDevice(ctx, room, id, uiForDev); err != nil {
		log.L.Warnf("unable to cache ui for device %s %s: %s", room, id, err)
	}

	return uiForDev, nil
}

func (d *dataService) cacheUIForDevice(ctx context.Context, room, id, uiForDev string) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_uiForDeviceBucket))
		if b == nil {
			return fmt.Errorf("ui for device bucket does not exist")
		}

		bytes, err := json.Marshal(uiForDev)
		if err != nil {
			return fmt.Errorf("unable to marshal ui for device: %w", err)
		}

		if err = b.Put([]byte(room+id), bytes); err != nil {
			return fmt.Errorf("unable to put ui for device: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *dataService) uiForDeviceFromCache(ctx context.Context, room, id string) (string, error) {
	var uiForDev string

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_uiForDeviceBucket))
		if b == nil {
			return fmt.Errorf("ui for device bucket does not exist")
		}

		bytes := b.Get([]byte(room + id))
		if bytes == nil {
			return fmt.Errorf("ui for device not in cache")
		}

		if err := json.Unmarshal(bytes, &uiForDev); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return uiForDev, err
	}

	return uiForDev, nil
}

func (d *dataService) ControlGroup(ctx context.Context, room, id string) (string, error) {
	controlGroup, err := d.dataService.ControlGroup(ctx, room, id)
	if err != nil {
		controlGroup, cacheErr := d.controlGroupFromCache(ctx, room, id)
		if cacheErr != nil {
			log.L.Warnf("unable to get control group for %s %s from cache: %s", room, id, cacheErr)
			return controlGroup, err
		}

		return controlGroup, nil
	}

	if err := d.cacheControlGroup(ctx, room, id, controlGroup); err != nil {
		log.L.Warnf("unable to cache controlGroup %s: %s", room, err)
	}

	return controlGroup, nil
}

func (d *dataService) cacheControlGroup(ctx context.Context, room, id, controlGroup string) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_controlGroupBucket))
		if b == nil {
			return fmt.Errorf("controlGroup bucket does not exist")
		}

		bytes, err := json.Marshal(controlGroup)
		if err != nil {
			return fmt.Errorf("unable to marshal controlGroup: %w", err)
		}

		if err = b.Put([]byte(room+id), bytes); err != nil {
			return fmt.Errorf("unable to put controlGroup: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *dataService) controlGroupFromCache(ctx context.Context, room, id string) (string, error) {
	var controlGroup string

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_controlGroupBucket))
		if b == nil {
			return fmt.Errorf("controlGroup bucket does not exist")
		}

		bytes := b.Get([]byte(room + id))
		if bytes == nil {
			return fmt.Errorf("controlGroup not in cache")
		}

		if err := json.Unmarshal(bytes, &controlGroup); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return controlGroup, err
	}

	return controlGroup, nil
}

func (d *dataService) RoomAndControlGroup(ctx context.Context, key string) (string, string, error) {
	room, controlGroup, err := d.dataService.RoomAndControlGroup(ctx, key)
	if err != nil {
		room, controlGroup, cacheErr := d.roomAndControlGroupFromCache(ctx, key)
		if cacheErr != nil {
			log.L.Warnf("unable to get room and control group %s from cache: %s", key, cacheErr)
			return room, controlGroup, err
		}

		return room, controlGroup, nil
	}

	forCache := roomAndControlgroup{
		Room:         room,
		ControlGroup: controlGroup,
	}

	if err := d.cacheRoomAndControlGroup(ctx, key, forCache); err != nil {
		log.L.Warnf("unable to cache room and control group %s: %s", key, err)
	}

	return room, controlGroup, nil
}

func (d *dataService) cacheRoomAndControlGroup(ctx context.Context, key string, forCache roomAndControlgroup) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_roomAndControlGroupBucket))
		if b == nil {
			return fmt.Errorf("room and control group bucket does not exist")
		}

		bytes, err := json.Marshal(&forCache)
		if err != nil {
			return fmt.Errorf("unable to marshal room and control group: %w", err)
		}

		if err = b.Put([]byte(key), bytes); err != nil {
			return fmt.Errorf("unable to put room and control group: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *dataService) roomAndControlGroupFromCache(ctx context.Context, key string) (string, string, error) {
	var fromCache roomAndControlgroup

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_roomAndControlGroupBucket))
		if b == nil {
			return fmt.Errorf("room and control group bucket does not exist")
		}

		bytes := b.Get([]byte(key))
		if bytes == nil {
			return fmt.Errorf("room and control group not in cache")
		}

		if err := json.Unmarshal(bytes, &fromCache); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return fromCache.Room, fromCache.ControlGroup, nil
}

func (d *dataService) Config(ctx context.Context, room string) (ui.Config, error) {
	config, err := d.dataService.Config(ctx, room)
	if err != nil {
		config, cacheErr := d.configFromCache(ctx, room)
		if cacheErr != nil {
			log.L.Warnf("unable to get config %s from cache: %s", room, cacheErr)
			return config, err
		}

		return config, nil
	}

	if err := d.cacheConfig(ctx, room, config); err != nil {
		log.L.Warnf("unable to cache config %s: %s", room, err)
	}

	return config, nil
}

func (d *dataService) cacheConfig(ctx context.Context, room string, config ui.Config) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_configBucket))
		if b == nil {
			return fmt.Errorf("config bucket does not exist")
		}

		bytes, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("unable to marshal config: %w", err)
		}

		if err = b.Put([]byte(room), bytes); err != nil {
			return fmt.Errorf("unable to put config: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *dataService) configFromCache(ctx context.Context, room string) (ui.Config, error) {
	var config ui.Config

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_configBucket))
		if b == nil {
			return fmt.Errorf("config bucket does not exist")
		}

		bytes := b.Get([]byte(room))
		if bytes == nil {
			return fmt.Errorf("config not in cache")
		}

		if err := json.Unmarshal(bytes, &config); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return config, err
	}

	return config, nil
}
