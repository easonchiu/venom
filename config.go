package venom

type Mode int

type Config struct {
	Mode          string
	Address       string
	Port          string
	SuccessFormat func(obj interface{}) interface{}
	FailFormat    func(errCode interface{}, errMessage string, obj interface{}) interface{}
}
