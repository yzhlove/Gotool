package service

type Service interface {
	Init() error
	Start() error
	Stop() error
}
