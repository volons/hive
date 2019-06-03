package models

import (
	"math"
	"time"
)

// Position contains positional data of a vehicle
type Position struct {
	Lat       float64   `json:"lat"`    // Latitude in degrees * 1E7
	Lon       float64   `json:"lon"`    // Longitude in degrees * 1E7
	Alt       float64   `json:"alt"`    // Altitude in millimeters
	RelAlt    float64   `json:"relAlt"` // Altitude above ground in millimeters
	Vx        float64   `json:"vx"`     // Ground X Speed (Latitude, positive north) in m/s * 100
	Vy        float64   `json:"vy"`     // Ground Y Speed (Longitude, positive east) in m/s * 100
	Vz        float64   `json:"vz"`     // Ground Z Speed (Altitude, positive down) in m/s * 100
	Hdg       float64   `json:"hdg"`    // Vehicle heading (yaw angle) in degrees * 100. If unknown, set to: UINT16_MAX
	Timestamp time.Time `json:"timestamp"`
}

// NewPosition creates and returns a new position model
func NewPosition(lat, lon, alt, relAlt, vx, vy, vz, hdg float64) Position {
	return Position{
		lat,
		lon,
		alt,
		relAlt,
		vx,
		vy,
		vz,
		hdg,
		time.Now(),
	}
}

// NewPoint creates a position with only lat, lon and relAlt values
func NewPoint(lat, lon, relAlt float64) Position {
	return Position{Lat: lat, Lon: lon, RelAlt: relAlt, Timestamp: time.Now()}
}

// LatInt returns the position's latitude value as an int32 in deg * 1e7
func (pos Position) LatInt() int32 {
	return int32(pos.Lat * 1e7)
}

// SetLat sets the position's latitude value in deg
func (pos *Position) SetLat(val float64) {
	pos.Lat = val
}

// LonInt returns the position's latitude value as an int32 in deg * 1e7
func (pos Position) LonInt() int32 {
	return int32(pos.Lon * 1e7)
}

// SetLon sets the position's longitude value in deg
func (pos *Position) SetLon(val float64) {
	pos.Lon = val
}

// SetAlt sets the position's altitude value in m
func (pos *Position) SetAlt(val float64) {
	pos.Alt = val
}

// SetRelAlt sets the position's relative altitude to the ground in m
func (pos *Position) SetRelAlt(val float64) {
	pos.RelAlt = val
}

// SetVx sets the position's latitudal speed in m/s
func (pos *Position) SetVx(val float64) {
	pos.Vx = val
}

// SetVy sets the position's longitudal speed in m/s
func (pos *Position) SetVy(val float64) {
	pos.Vy = val
}

// SetVz sets the position's vertical speed in m/s
func (pos *Position) SetVz(val float64) {
	pos.Vz = val
}

// SetHdg sets the position's heading value in deg
func (pos *Position) SetHdg(val float64) {
	pos.Hdg = val
}

// Translate returns a point translated by x, y, z meters
func (pos Position) Translate(x float64, y float64, z float64) Position {
	out := pos
	if x != 0 || y != 0 {
		rEarth := float64(6371000)
		rad := math.Pi / 180
		deg := 180 / math.Pi

		out.Lat = pos.Lat + ((y / rEarth) * deg)
		latRad := out.Lat * rad
		out.Lon = pos.Lon + ((x / rEarth) * deg / math.Cos(latRad))
	}

	if z != 0 {
		out.Alt = pos.Alt + z
		out.RelAlt = pos.RelAlt + z
	}

	return out
}

// Diff calculates the diff between to position
func (pos Position) Diff(p2 Position) (float64, float64, float64) {
	dx := p2.Lon - pos.Lon
	dy := p2.Lat - pos.Lat
	dz := p2.Alt - pos.Alt

	rEarth := float64(6371000)
	rad := math.Pi / 180
	//deg := 180 / math.Pi

	var x, y float64

	if dx != 0 {
		y = dy * rad * rEarth
	}

	if dy != 0 {
		latRad := pos.Lat * rad
		x = dx * math.Cos(latRad) * rad * rEarth
	}

	return x, y, dz
}

// Distance calculates the distance between to position
func (pos Position) Distance(p2 Position) float64 {
	x, y, z := pos.Diff(p2)
	return math.Sqrt(math.Sqrt(x*x+y*y) + z*z)
}
