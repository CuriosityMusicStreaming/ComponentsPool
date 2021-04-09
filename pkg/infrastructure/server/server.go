package server

type StopFunc func() error

type Server interface {
	Serve() error
	Stop() error
}

type FuncServer struct {
	ServeImpl func() error
	StopImpl  func() error
}

func (f *FuncServer) Serve() error {
	return f.ServeImpl()
}

func (f *FuncServer) Stop() error {
	return f.StopImpl()
}
