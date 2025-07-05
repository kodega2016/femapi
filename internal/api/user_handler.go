package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/kodega2016/femapi/internal/store"
	"github.com/kodega2016/femapi/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New(`email address is required`)
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Panicf("ERROR: decoding handle register request: %w", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "failed to decode the request",
		})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		h.logger.Printf("ERROR: validateRegisterRequest %w", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": err.Error(),
		})
		return
	}
	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
}
