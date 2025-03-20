package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	"berezovskyivalerii/csv-rest-app/internal/domain"
	"encoding/csv"
)

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		logError("create", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := h.productService.Create(ctx, product)
	if err != nil {
		logError("create", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		logError("create", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	products, err := h.productService.GetAll(ctx)
	if err != nil {
		logError("getAll", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check format request (JSON или CSV)
	format := r.URL.Query().Get("format")

	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=products.csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		err := writer.Write([]string{"Name", "Price"})
		if err != nil{
			logError("getAll", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, product := range products {
			writer.Write([]string{
				product.Name,
				strconv.FormatInt(product.Price, 10),
			})
		}
	} else {
		// JSON used by default
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		logError("update", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product domain.ProductUpdate
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		logError("update", err)
		return
	}

	updatedProductID, err := h.productService.Update(r.Context(), id, product)
	if err != nil {
		logError("update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedProductID); err != nil {
		logError("update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		logError("delete", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deletedProductID, err := h.productService.Delete(r.Context(), id)
	if err != nil {
		logError("delete", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(deletedProductID); err != nil {
		logError("delete", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
