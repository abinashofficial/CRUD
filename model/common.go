package model

var StudentInfo = []TestCase{{Name: "Abi", MobNumber: 994048389}, {Name: "Arun", MobNumber: 987436652}}

type TestCase struct {
	Name      string
	MobNumber int
}

type Requests struct {
	Name string `json:"Name"`
}
