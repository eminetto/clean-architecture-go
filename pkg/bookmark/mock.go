package bookmark

//NewServiceMock return a new mock service to be used in tests
func NewServiceMock() *Service {
	repo := NewInmemRepository()
	return NewService(repo)
}
