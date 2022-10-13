package models

type Thermometer struct {
	Raised float64 `json:"raised"`
	Goal   uint64  `json:"goal"`
}

type AllThermometers map[string]*Thermometer
