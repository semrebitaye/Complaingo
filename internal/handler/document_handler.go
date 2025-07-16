package handler

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/middleware"
	"Complaingo/internal/usecase"
	"Complaingo/internal/utility"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	appErrors "Complaingo/internal/errors"

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

	middleware.WriteSuccess(w, doc, "File uploaded successfully", http.StatusCreated)
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

	query := r.URL.Query()

	paginParam := utility.PaginationParam{
		Page:    query.Get("page"),
		PerPage: query.Get("per_page"),
		Sort:    query.Get("sort"),
		Search:  query.Get("search"),
		Filter:  query.Get("filter"),
	}

	filterParam, err := utility.ExtractPagination(paginParam)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.Wrap(err, "Failed to prase query"))
		return
	}

	doc, err := h.usecase.GetDocumentByUser(r.Context(), userID, filterParam)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteSuccess(w, doc, "files retrieved successfully by user id", http.StatusOK)
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

	middleware.WriteSuccess(w, "", "Document deleted successfully", http.StatusNoContent)
}
