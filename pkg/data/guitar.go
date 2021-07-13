package data

type Guitar struct {
	Brand string
	Model string
}

type FinishColor struct {
	Color string `json:"color"`
}

var FinishColors = []FinishColor{
	{"Sunburst"},
	{"Natural"},
	{"Red"},
}

var Guitars = map[string]string{
	"Gibson": "Les Paul",
	"Fender": "Telecaster",
}

var Amps = map[string]string{
	"Fender": "Deluxe",
	"Vox":    "AC30",
}
