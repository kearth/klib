package server

type Server struct{}

func (s *Server) Run() error {
	return nil
}

func (s *Server) Shutdown() error {
	return nil
}

func NewServer() *Server {
	return &Server{}
}
