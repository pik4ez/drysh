package out

type Out interface {
	Write(data string) string
	GetConfig() interface{}
}
