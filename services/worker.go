package services

type Worker interface {
	Start()
	Stop()
	Name() string
}
