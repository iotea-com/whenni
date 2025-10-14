package things

type Attributes interface {
	Category() ThingCategory
	Validate() error
	MarshalJson() ([]byte, error)
}
