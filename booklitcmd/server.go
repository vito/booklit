package booklitcmd

import (
	"net/http"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Command    *Command
	FileServer http.Handler

	buildLock sync.Mutex
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Debugln("serving", r.URL.Path)

	if strings.HasSuffix(r.URL.Path, ".html") {
		logrus.Info("building")

		server.buildLock.Lock()

		err := server.Command.Build()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			server.buildLock.Unlock()
			return
		}

		logrus.Info("build complete")

		server.buildLock.Unlock()
	}

	server.FileServer.ServeHTTP(w, r)
}
