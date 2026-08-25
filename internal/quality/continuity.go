package quality

import "github.com/jb843051627/fjord-resonance/internal/model"

func Continuity(samples []model.AcousticSample) float64 {
	if len(samples) == 0 {
		return 0
	}
	if len(samples) == 1 {
		return 1
	}
	ordered := 0
	for i := 1; i < len(samples); i++ {
		if samples[i].Sequence == samples[i-1].Sequence+1 {
			ordered++
		}
	}
	return float64(ordered) / float64(len(samples)-1)
}

func HasDuplicateSequence(samples []model.AcousticSample) bool {
	seen := make(map[int]struct{}, len(samples))
	for _, sample := range samples {
		if _, ok := seen[sample.Sequence]; ok {
			return true
		}
		seen[sample.Sequence] = struct{}{}
	}
	return false
}

func MissingSequences(samples []model.AcousticSample) []int {
	if len(samples) < 2 {
		return nil
	}
	min, max := samples[0].Sequence, samples[0].Sequence
	seen := make(map[int]bool, len(samples))
	for _, sample := range samples {
		if sample.Sequence < min {
			min = sample.Sequence
		}
		if sample.Sequence > max {
			max = sample.Sequence
		}
		seen[sample.Sequence] = true
	}
	missing := make([]int, 0)
	for value := min; value <= max; value++ {
		if !seen[value] {
			missing = append(missing, value)
		}
	}
	return missing
}
