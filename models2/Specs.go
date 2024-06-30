package models

type Specs struct {
	LicensePlate string `json:"licensePlate" gorm:"primaryKey"`
	Brand        string `json:"brand"`
	Color        string `json:"color,omitempty"`
	EngineSize   int    `json:"engineSize,omitempty"`
	FirstReg     string `json:"firstReg,omitempty"`
	FirstRegHun  string `json:"firstRegHun,omitempty"`
	FuelType     string `json:"fuelType,omitempty"`
	Gearbox      string `json:"gearbox,omitempty"`
	Model        string `json:"model"`
	NumOfOwners  int    `json:"numOfOwners,omitempty"`
	Performance  int    `json:"performance,omitempty"`
	Status       string `json:"status,omitempty"`
	TypeCode     string `json:"typeCode"`
	Year         int    `json:"year"`
}
