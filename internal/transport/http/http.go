package httptransport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/k1ender/psf/internal/middleware"
	"github.com/k1ender/psf/internal/repository"
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
		ctx := r.Context()
		log := middleware.FromContext(ctx)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data, err := templates.GetTemplate(templates.Upload)
		if err != nil {
			log.ErrorContext(ctx, "failed to get template", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write(data)
	})

	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := middleware.FromContext(ctx)

		file, headers, err := r.FormFile("file")
		if err != nil {
			log.ErrorContext(ctx, "failed to get file", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		id, err := s.fileService.SaveFile(ctx, file, headers.Filename)
		if err != nil {
			log.ErrorContext(ctx, "failed to save file", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		defer file.Close()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, id)
	})

	mux.HandleFunc("GET /file/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := middleware.FromContext(ctx)

		id := r.PathValue("id")

		file, filename, err := s.fileService.GetFile(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				log.InfoContext(ctx, "file not found")
				http.NotFound(w, r)
				return
			}

			log.ErrorContext(ctx, "failed to get file", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		extension := filepath.Ext(filename)
		extension = strings.TrimPrefix(extension, ".")

		w.Header().Set("Content-Type", mime.TypeByExtension("."+extension))
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		w.Header().Set("Content-Length", fmt.Sprint(len(file)))

		_, err = w.Write(file)
		if err != nil {
			log.ErrorContext(ctx, "failed to write file", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	s.RegisterRoutes(mux)

	var handler http.Handler = mux
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)

	s.httpServer.Handler = handler

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil && err != http.ErrServerClosed {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
