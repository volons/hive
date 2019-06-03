package models

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"os"
	"time"
)

// The Fence class allows to restrain a vehicles movements
type Fence struct {
	a Position
	b Position
}

// PointData represents a point in the world
type PointData struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt"`
}

// FenceData is a json serializable version of fence struct
type FenceData []PointData

var fenceInstance *Fence

// GetFence returns the current fence
func GetFence() *Fence {
	return fenceInstance
}

// SetFence sets the current fence
func SetFence(data FenceData) error {
	if len(data) != 2 {
		return errors.New("Fence should have exactly two points")
	}

	fence := &Fence{
		a: NewPoint(data[0].Lat, data[0].Lon, data[0].Alt),
		b: NewPoint(data[1].Lat, data[1].Lon, data[1].Alt),
	}

	if !fence.Valid() {
		return errors.New("fence is not valid")
	}

	fenceInstance = fence

	return nil
}

// SetFenceJSON sets the current fence in json format
func SetFenceJSON(data []byte) error {
	var fence FenceData
	err := json.Unmarshal(data, &fence)
	if err != nil {
		return err
	}

	return SetFence(fence)
}

// SetFenceFile creates the current fence according to the specified file
func SetFenceFile(fenceFilePath string) error {
	fenceFile, err := os.Open(fenceFilePath)
	if err != nil {
		return err
	}

	var data FenceData
	decoder := json.NewDecoder(fenceFile)
	if err = decoder.Decode(&data); err != nil {
		return err
	}

	return SetFence(data)
}

// A return the a point
func (fence *Fence) A() Position {
	return fence.a
}

// B return the b point
func (fence *Fence) B() Position {
	return fence.b
}

// Check checks if the provided position is inside the fence
func (fence *Fence) Check(pos Position) bool {
	if pos.Lon > fence.b.Lon ||
		pos.Lon < fence.a.Lon ||
		pos.Lat > fence.b.Lat ||
		pos.Lat < fence.a.Lat ||
		pos.RelAlt > fence.b.RelAlt ||
		pos.RelAlt < fence.a.RelAlt {
		return false
	}

	return true
}

// GetAutoRc returns the corrcted rc values to not leave the fence
func (fence *Fence) GetAutoRc(pos Position, manualRc *Rc) (*Rc, *Position, bool, bool) {
	var autoRc *Rc
	var slowed, outside bool
	var target *Position

	if !fence.Check(pos) {
		outside = true
		target = &Position{}
		*target = pos

		var tX, tY, tZ float64

		if target.Lon > fence.b.Lon {
			target.SetLon(fence.b.Lon)
			tX = -5
		} else if target.Lon < fence.a.Lon {
			target.SetLon(fence.a.Lon)
			tX = 5
		}

		if target.Lat > fence.b.Lat {
			target.SetLat(fence.b.Lat)
			tY = -5
		} else if target.Lat < fence.a.Lat {
			target.SetLat(fence.a.Lat)
			tY = 5
		}

		if target.RelAlt > fence.b.RelAlt {
			target.SetRelAlt(fence.b.RelAlt)
			tZ = -5
		} else if target.RelAlt < fence.a.RelAlt {
			target.SetRelAlt(fence.a.RelAlt)
			tZ = 5
		}

		*target = target.Translate(tX, tY, tZ)
		autoRc = NewNullRc()
	} else {
		autoRc = manualRc.Copy()

		m := 15.0
		f := pos.Translate(0, m, 0)
		b := pos.Translate(0, -m, 0)
		l := pos.Translate(-m, 0, 0)
		r := pos.Translate(m, 0, 0)
		u := pos.Translate(0, 0, m)
		d := pos.Translate(0, 0, -m)

		xa, ya, za := fence.meterDist(pos, fence.a)
		xb, yb, zb := fence.meterDist(pos, fence.b)

		limits := NewRcLimits()

		if !fence.Check(r) {
			limits.RollMax = xb / m
		}
		if !fence.Check(l) {
			limits.RollMin = xa / m
		}
		if !fence.Check(f) {
			limits.PitchMax = yb / m
		}
		if !fence.Check(b) {
			limits.PitchMin = ya / m
		}
		if !fence.Check(u) {
			limits.ThrottleMax = zb / m
		}
		if !fence.Check(d) {
			limits.ThrottleMin = za / m
		}

		autoRc.Rotate(-pos.Hdg)
		slowed = autoRc.ApplyLimits(limits)
		autoRc.Rotate(pos.Hdg)

		//debug(pos, limits, autoRc)
	}

	return autoRc, target, slowed, outside
}

var last time.Time

func debug(pos Position, limits RcLimits, rc *Rc) {
	if time.Since(last) > time.Millisecond*500 {
		log.Printf("pos lat: %v, lon: %v, relAlt: %v, hdg: %v", pos.Lat, pos.Lon, pos.RelAlt, pos.Hdg)
		log.Printf("roll min: %v, max: %v", limits.RollMin, limits.RollMax)
		log.Printf("pitch min: %v, max: %v", limits.PitchMin, limits.PitchMax)
		log.Printf("throttle min: %v, max: %v", limits.ThrottleMin, limits.ThrottleMax)
		log.Printf("rc roll: %v, pitch: %v, throttle: %v", rc.Roll(), rc.Pitch(), rc.Throttle())
		last = time.Now()
	}
}

func (fence *Fence) meterDist(p1 Position, p2 Position) (x, y, z float64) {
	rEarth := float64(6371000)
	//PI := 3.1415926535
	rad := math.Pi / 180

	deg := 180 / math.Pi
	diffLat := p2.Lat - p1.Lat
	diffLon := p2.Lon - p1.Lon
	diffAlt := p2.RelAlt - p1.RelAlt
	latRad := p1.Lat * rad

	return (diffLon * math.Cos(latRad) / deg) * rEarth, (diffLat / deg) * rEarth, diffAlt
}

// PrintDiff prints the distance between the fences
// points and the supplied position
func (fence *Fence) PrintDiff(pos Position) {
	log.Printf("Fence diff: Lat: %v, Lon: %v\n", pos.Lat-fence.a.Lat, pos.Lon-fence.a.Lon)
	log.Printf("Fence diff: Lat: %v, Lon: %v\n", pos.Lat-fence.b.Lat, pos.Lon-fence.b.Lon)
}

// JSON returns a to JSON convertable representation of this fence
func (fence *Fence) JSON() FenceData {
	return FenceData{
		PointData{fence.A().Lat, fence.A().Lon, fence.A().RelAlt},
		PointData{fence.B().Lat, fence.B().Lon, fence.B().RelAlt},
	}
}

// Valid checks if the fence is valid
func (fence *Fence) Valid() bool {
	return fence.A().Lat < fence.B().Lat &&
		fence.A().Lon < fence.B().Lon &&
		fence.A().RelAlt < fence.B().RelAlt
}
