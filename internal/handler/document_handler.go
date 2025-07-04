package handler

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/middleware"
	"Complaingo/internal/usecase"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

type DocumentHandler struct {
	usecase *usecase.DocumentUsecase
}

func NewDocumentHandler(usecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{
		usecase: usecase,
	}
}

func (h *DocumentHandler) Uplod(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserId(r.Context())

	// parse uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	defer file.Close()

	// save to disc
	destPath := filepath.Join("uploads_doc", header.Filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	defer destFile.Close()
	io.Copy(destFile, file)

	doc := &models.Document{
		UserID:   userID,
		FileName: header.Filename,
		FilePath: destPath,
	}

	if err := h.usecase.Uplod(r.Context(), doc); err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteSuccess(w, doc, "File uploaded successfully")
}

func (h *DocumentHandler) GetDocumentByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	doc, err := h.usecase.GetDocumentByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	http.ServeFile(w, r, doc.FilePath)
}

func (h *DocumentHandler) GetDocumentByUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	doc, err := h.usecase.GetDocumentByUser(r.Context(), userID)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteSuccess(w, doc, "files retrieved successfully by user id")
}

func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	if err := h.usecase.DeleteDocument(r.Context(), id); err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteSuccess(w, "", "Document deleted successfully")
}
