// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Translation -.
type Camera struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Firmware     string `json:"firmware"`
	Serial       string `json:"serial"`
	IP           string `json:"ip"`
}
