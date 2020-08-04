package param

type (
	// Param path parameter
	Param struct {
		Key   string
		Value string
	}
	// Params path parameters
	Params []Param

	// IParam path parameters interface
	IParam interface {
		ByName(string) string
	}
)

// ByName  get paramter name
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}
