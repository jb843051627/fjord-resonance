package service

import "github.com/jb843051627/fjord-resonance/internal/store"

func structToBuoyFilter() store.BuoyFilter {
	filter := store.BuoyFilter{Limit: 200}
	return filter
}

func structToSampleFilter() store.SampleFilter {
	return store.SampleFilter{Limit: 10000, OnlyValid: true}
}
