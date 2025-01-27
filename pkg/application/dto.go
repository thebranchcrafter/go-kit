package application

type Dto interface {
	Id() string
}

type InvalidDto struct {
	message string
}

func (i InvalidDto) Error() string {
	return i.message
}

type Command interface {
	Dto
}

type Query interface {
	Dto
}
