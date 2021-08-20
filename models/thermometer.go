package models

type Thermometer struct {
	Raised float64 `json:"raised"`
	Goal   uint64  `json:"goal"`
}

type AllThermometers struct {
	Halifax   *Thermometer `json:"halifax"`
	Calgary   *Thermometer `json:"calgary"`
	Edmonton  *Thermometer `json:"edmonton"`
	Montreal  *Thermometer `json:"montreal"`
	StJohns   *Thermometer `json:"stjohns"`
	Toronto   *Thermometer `json:"toronto"`
	Vancouver *Thermometer `json:"vancouver"`
}
