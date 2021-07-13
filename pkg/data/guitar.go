package data

// Guitar struct used for marshalling/unmarshalling json key/value pairs
type Guitar struct {
	Brand string `json:"brand"` // this attribute is not required, lets Go know to rename the field "brand", as exported (public) fields requires uppercase
	Model string `json:"model"`

	// In Go, you can embed types and then have access to Guitar.Color or Guitar.Finish.Color
	// It also means a guitar gets any method a finish would have
	Finish
}

type Finish struct {
	Color string `json:"color"`
}

// Example to show how embedded structs enable their types to be used
// This section isn't used by the APIs, but just for illustration
func (f Finish) GetColor() string {
	return f.Color
}

func GuitarGetColor(g Guitar) string {
	return g.GetColor()
}
