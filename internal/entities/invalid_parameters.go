package entities

type InvalidParameters struct {
	Params []string
}

func NewInvalidParams() *InvalidParameters {
	return &InvalidParameters{}
}
func (i *InvalidParameters) Error() string {
	return "Invalid parameters"
}
func (i *InvalidParameters) AddParameters(value []string) {
	i.Params = append(i.Params, value...)
}
