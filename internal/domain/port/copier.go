package port

type Copier interface {
	Copy(dst, src interface{}) error
}
