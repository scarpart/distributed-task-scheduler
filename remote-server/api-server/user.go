package apiserver

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserAPIKeyParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) CreateUserAPIKey(ctx *gin.Context) {
	var req CreateUserAPIKeyParams 
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}
	
	// Creating the transaction to be committed, ensuring rollbacks in case it fails.
	tx, err := server.store.BeginTx(ctx, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// The rollback is ignored if there is a commit later on
	defer tx.Rollback()

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, genericErrorResponse("Invalid username provided."))
		return 
	}

	// Need to compare the bytes instead of strings in order to avoid errors with null bytes and such
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, genericErrorResponse("Invalid username or password provided."))
		return
	}

	apiKey, err := generateAPIKeys(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	arg := db.SetUserAPIKeyParams{
		ApiKey: sql.NullString{
			String: string(apiKey[:]),
			Valid: true,
		},
		UserID: user.UserID,
	}

	err = server.store.SetUserAPIKey(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	// Ends the transaction on a commit. If there is an error, we rollback
	if err = tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Header("API-Key", apiKey)
	ctx.JSON(http.StatusOK, "API Key successfully created.")
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"    binding:"required"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := hashPassword(req.Password)	
	if err != nil {
		// Probably change this error message sometime soon
		ctx.JSON(http.StatusInternalServerError, genericErrorResponse("An error occurred, please contact an admin."))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)	
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}

func generateAPIKeys(length int32) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", nil
	}
	apiKey := base64.URLEncoding.EncodeToString(randomBytes)
	return apiKey, nil
}

