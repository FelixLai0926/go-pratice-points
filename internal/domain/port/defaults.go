package port

type Defaults interface {
	Set(ptr interface{}) error
}
