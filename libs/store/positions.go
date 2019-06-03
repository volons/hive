package store

import (
	"sync"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/models"
)

type positionList struct {
	positions *sync.Map
}

func newPositionList() positionList {
	return positionList{
		positions: &sync.Map{},
	}
}

// Add sets a position of the position list
func (v positionList) Set(id string, position models.Position) {
	var p *models.Position
	p = &position
	v.positions.Store(key(id), p)
}

// Remove removes a position from the position list
func (v positionList) Remove(id string) {
	v.positions.Delete(key(id))
}

// Get returns a position by vehicle ID
func (v positionList) Get(id string) *models.Position {
	val, ok := v.positions.Load(key(id))
	if !ok {
		return nil
	}
	position, ok := val.(*models.Position)
	if !ok {
		return nil
	}

	return position
}

func (v positionList) Length() int {
	len := 0
	v.positions.Range(func(key interface{}, val interface{}) bool {
		len++
		return true
	})

	return len
}

// JSON returns the list of positions in a json serializable format
func (v positionList) JSON(since time.Time) (libs.JSONObject, int) {
	list := libs.JSONObject{}
	len := 0

	v.positions.Range(func(key interface{}, val interface{}) bool {
		v, ok := val.(*models.Position)
		if ok /*&& v.Timestamp().After(since)*/ {
			len++
			list[key.(string)] = v
		}
		return true
	})

	return list, len
}
