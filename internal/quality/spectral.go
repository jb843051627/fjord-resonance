package quality

import (
	"math"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type SpectralBand struct {
	LowHz  float64
	HighHz float64
	Weight float64
}

func EnergyByBand(samples []model.AcousticSample, bands []SpectralBand) []float64 {
	energy := make([]float64, len(bands))
	for _, sample := range samples {
		if !sample.Valid {
			continue
		}
		for index, band := range bands {
			if sample.FrequencyHz >= band.LowHz && sample.FrequencyHz < band.HighHz {
				energy[index] += math.Abs(sample.AmplitudeDB-sample.NoiseDB) * band.Weight
				break
			}
		}
	}
	return energy
}

func DominantBand(samples []model.AcousticSample, bands []SpectralBand) (SpectralBand, bool) {
	energy := EnergyByBand(samples, bands)
	best := -1
	for index, value := range energy {
		if best < 0 || value > energy[best] {
			best = index
		}
	}
	if best < 0 || best >= len(bands) {
		return SpectralBand{}, false
	}
	return bands[best], energy[best] > 0
}

func BandBalanced(energy []float64, tolerance float64) bool {
	if len(energy) < 2 {
		return true
	}
	mean := Mean(energy)
	for _, value := range energy {
		if math.Abs(value-mean) > tolerance {
			return false
		}
	}
	return true
}
