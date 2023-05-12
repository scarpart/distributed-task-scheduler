package apiserver

import (
	"crypto/sha256"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Creating middleware functions to validate API Keys
func (server *Server) ValidateKeys() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("API-Key")
		hashedKey := sha256.Sum256([]byte(apiKey))
	
		// TODO: change this to :exec (no need for a return)
		_, err := server.store.GetAPIKeys(ctx, sql.NullString{
			String: string(hashedKey[:]),
			Valid: true,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	
		ctx.Next()
	}
}

