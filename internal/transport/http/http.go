package httptransport

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/k1ender/psf/internal/service"
	"github.com/k1ender/psf/templates"
)

type Server struct {
	httpServer  *http.Server
	fileService service.File
}

func New(addr string, fileService service.File) *Server {
	server := &http.Server{
		Addr: addr,
	}

	return &Server{
		httpServer:  server,
		fileService: fileService,
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data, err := templates.GetTemplate(templates.Upload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	})

	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		file, headers, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := s.fileService.SaveFile(file, headers.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, id)
	})

	mux.HandleFunc("GET /file/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		file, filename, err := s.fileService.GetFile(id)
		extension := filepath.Ext(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ext := strings.TrimPrefix(extension, ".")
		w.Header().Set("Content-Type", mime.TypeByExtension("."+ext))
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	s.RegisterRoutes(mux)

	s.httpServer.Handler = mux

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err == http.ErrServerClosed {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
