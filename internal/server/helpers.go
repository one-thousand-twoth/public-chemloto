package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
)

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func encodeError(w http.ResponseWriter, log *slog.Logger, err error) {
	enerr.HTTPErrorResponse(w, log, err)
}
