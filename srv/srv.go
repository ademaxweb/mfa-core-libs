package srv

import (
	"fmt"
	"github.com/ademaxweb/mfa-core-libs/discovery"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	Router    *mux.Router
	S         *http.Server
	Discovery *discovery.Discovery
}

func New(port int, timeout time.Duration, serviceDiscovery *discovery.Discovery) *Server {
	router := mux.NewRouter()

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Handler:      router,
	}

	return &Server{
		Router:    router,
		S:         s,
		Discovery: serviceDiscovery,
	}
}

func (s *Server) Run() error {
	err := s.Discovery.Register()
	if err != nil {
		return err
	}
	return s.S.ListenAndServe()
}
