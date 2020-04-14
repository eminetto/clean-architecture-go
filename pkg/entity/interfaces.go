package entity

type Validable interface {
	Validate() error
}
