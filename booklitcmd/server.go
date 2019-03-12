package booklitcmd

import (
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/vito/booklit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Server struct {
	In        string
	Processor *load.Processor

	Templates string
	Engine    *render.HTMLRenderingEngine

	FileServer http.Handler

	buildLock sync.Mutex
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"request": r.URL.Path,
	})

	log.Debugln("serving")

	section, found, err := server.loadRequestedSection(r.URL.Path)
	if err != nil {
		log.Errorf("failed to load section: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		server.FileServer.ServeHTTP(w, r)
		return
	}

	server.buildLock.Lock()
	defer server.buildLock.Unlock()

	if server.Templates != "" {
		err := server.Engine.LoadTemplates(server.Templates)
		if err != nil {
			log.Errorf("failed to load templates: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	log = log.WithFields(logrus.Fields{
		"section": section.Path,
	})

	log.Info("rendering")

	err = server.Engine.RenderSection(w, section)
	if err != nil {
		log.Errorf("failed to render: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func (server *Server) loadRequestedSection(path string) (*booklit.Section, bool, error) {
	ext := server.Engine.FileExtension()

	if path == "/" {
		path = "/index." + ext
	}

	if !strings.HasSuffix(path, "."+ext) {
		return nil, false, nil
	}

	tagName := strings.TrimSuffix(filepath.Base(path), "."+ext)

	logrus.WithFields(logrus.Fields{
		"section": server.In,
	}).Info("loading root section")

	rootSection, err := server.Processor.LoadFile(server.In, basePluginFactories)
	if err != nil {
		return nil, false, err
	}

	tags := rootSection.FindTag(tagName)
	if len(tags) == 0 {
		return nil, false, nil
	}

	return tags[0].Section, true, nil
}
