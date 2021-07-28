package models

type Port struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	City        string      `json:"city"`
	Alias       []string    `json:"alias"`
	Regions     []string    `json:"regions"`
	Coordinates Coordinates `json:"coordinates"`
	Province    string      `json:"province"`
	Timezone    string      `json:"timezone"`
	Unlocs      []string    `json:"unlocs"`
}
