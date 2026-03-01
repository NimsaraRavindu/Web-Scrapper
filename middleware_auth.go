package main

import (
	"fmt"
	"net/http"

	"github.com/NimsaraRavindu/webScrapper/internal/auth"
	"github.com/NimsaraRavindu/webScrapper/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Failed to get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
