package model

import "math"

type Coordinate struct{ Latitude, Longitude float64 }

func (c Coordinate) Valid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 && c.Longitude >= -180 && c.Longitude <= 180
}

func (c Coordinate) DistanceTo(other Coordinate) float64 {
	const radius = 6371000.0
	lat1, lat2 := c.Latitude*math.Pi/180, other.Latitude*math.Pi/180
	dLat := (other.Latitude - c.Latitude) * math.Pi / 180
	dLon := (other.Longitude - c.Longitude) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return radius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (b Buoy) Coordinate() Coordinate {
	return Coordinate{Latitude: b.Latitude, Longitude: b.Longitude}
}

func WithinRadius(origin Coordinate, candidates []Buoy, radius float64) []Buoy {
	result := make([]Buoy, 0)
	for _, buoy := range candidates {
		if origin.DistanceTo(buoy.Coordinate()) <= radius {
			result = append(result, buoy)
		}
	}
	return result
}
