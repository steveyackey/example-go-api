package data

var Finishs = []Finish{
	{"Sunburst"},
	{"Natural"},
	{"Red"},
}

var Guitars = map[string]Guitar{
	"Gibson": {Brand: "Gibson", Model: "Les Paul", Finish: Finish{Color: "Sunburst"}},
	"Fender": {Brand: "Fender", Model: "Stratocaster", Finish: Finish{Color: "Sunburst"}},
}
