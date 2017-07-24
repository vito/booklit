package booklitcmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Command    *Command
	FileServer http.Handler

	lastBuilt  time.Time
	builtPaths []string
	buildLock  sync.Mutex
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Debugln("serving", r.URL.Path)

	server.buildLock.Lock()

	if server.shouldBuild() {
		paths, err := server.Command.Build(false)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build:\n\n%s", err), http.StatusInternalServerError)
			server.buildLock.Unlock()
			return
		}

		server.builtPaths = paths
		server.lastBuilt = time.Now()

		logrus.Info("build complete")
	}

	server.buildLock.Unlock()

	server.FileServer.ServeHTTP(w, r)
}

func (server *Server) shouldBuild() bool {
	if server.builtPaths == nil {
		logrus.Info("initial build")
		return true
	}

	wd, _ := os.Getwd()

	for _, path := range server.builtPaths {
		logPath := path
		if filepath.IsAbs(path) {
			relPath, err := filepath.Rel(wd, path)
			if err != nil {
				logrus.Errorf("failed to resolve relative path for %s: %s", path, err)
			} else {
				logPath = relPath
			}
		}

		log := logrus.WithFields(logrus.Fields{
			"path": logPath,
		})

		info, err := os.Stat(path)
		if err != nil {
			log.Infof("file removed; rebuilding")
			return true
		}

		if info.ModTime().After(server.lastBuilt) {
			log.Infof("change detected; rebuilding")
			return true
		}
	}

	return false
}
