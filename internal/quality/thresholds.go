package quality

type Thresholds struct {
	MinCoverage     float64
	MaxNoiseFloorDB float64
	MaxDriftDB      float64
	MinContinuity   float64
	ReviewScore     float64
	PassScore       float64
}

func DefaultThresholds() Thresholds {
	return Thresholds{MinCoverage: 0.70, MaxNoiseFloorDB: -32, MaxDriftDB: 5, MinContinuity: 0.90, ReviewScore: 0.65, PassScore: 0.82}
}

func (t Thresholds) Valid() bool {
	return t.MinCoverage > 0 && t.MaxNoiseFloorDB < 0 && t.MaxDriftDB > 0 && t.MinContinuity > 0 && t.PassScore > t.ReviewScore
}

func (t Thresholds) Score(coverage, noise, drift, continuity float64) float64 {
	coverageScore := clamp(coverage, 0, 1)
	noiseScore := clamp(1-(noise-t.MaxNoiseFloorDB)/20, 0, 1)
	driftScore := clamp(1-drift/t.MaxDriftDB, 0, 1)
	continuityScore := clamp(continuity, 0, 1)
	return coverageScore*0.35 + noiseScore*0.2 + driftScore*0.2 + continuityScore*0.25
}

func clamp(value, low, high float64) float64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
