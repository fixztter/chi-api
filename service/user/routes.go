package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fixztter/chi-api/config"
	"github.com/fixztter/chi-api/service/auth"
	"github.com/fixztter/chi-api/types"
	"github.com/fixztter/chi-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Post("/register", h.handleRegister)
	router.Post("/login", h.handleLogin)
	router.Get("/login/delete/{userID}", h.handleDelete)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	err := utils.ParseJSON(r, &user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unexpected error: %s", err.Error()))
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	err := utils.ParseJSON(r, &user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	_, err = h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %v already exists", user.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error parsing user_id: %v", err))
		return
	}

	if err := h.store.DeleteUser(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error traying to delete user: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully!"})
}
