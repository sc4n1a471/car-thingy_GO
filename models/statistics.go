package models

type Statistics struct {
	CarCount    int               `json:"carCount"`
	KnownCars   int               `json:"knownCars"`
	UnknownCars int               `json:"unknownCars"`
	BrandStats  []BrandStatistics `json:"brandStats"`
}

type BrandStatistics struct {
	Brand  string            `json:"brand"`
	Count  int               `json:"count"`
	Models []ModelStatistics `json:"models"`
}

type ModelStatistics struct {
	Model string `json:"model"`
	Count int    `json:"count"`
}
