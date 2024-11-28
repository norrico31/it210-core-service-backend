package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/entities"
	"github.com/norrico31/it210-core-service-backend/utils"
)

// TODO: REFACTO ALL OF THE CRUD HERE
type Handler struct {
	store entities.UserStore
}

func NewHandler(store entities.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	payload := entities.UserLoginPayload{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}

	token, user, err := h.store.Login(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":           user.ID,
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"age":          user.Age,
			"lastActiveAt": user.LastActiveAt,
			"createdAt":    user.CreatedAt,
			"updatedAt":    user.UpdatedAt,
			"deletedAt":    user.DeletedAt,
		},
		"token": token,
	})
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	for key, values := range r.Header {
		fmt.Printf("Header: %s, Value: %v\n", key, values)
	}

	users, err := h.store.GetUsers()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": users})
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user Id"))
		return
	}
	userId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}
	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": user})
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload entities.UserCreatePayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	password, err := utils.HashPassword("secret123")
	if err != nil {
		return
	}

	// FEATURE
	// SEND AN EMAIL FOR THE PASSWORD OF USER
	// ACTIVATE THE ACCOUNT (MAYBE?)
	err = h.store.CreateUser(entities.UserCreatePayload{
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		Email:      payload.Email,
		Age:        payload.Age,
		RoleId:     payload.RoleId,
		ProjectIDS: payload.ProjectIDS,
		Password:   password,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from URL
	vars := mux.Vars(r)
	str, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	// Parse payload
	var payload entities.UserUpdatePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Fetch existing user
	userExist, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch user: %v", err))
		return
	}

	// Prepare updated user
	user := entities.UserUpdatePayload{}
	if payload.FirstName != nil {
		user.FirstName = payload.FirstName
	} else {
		user.FirstName = &userExist.FirstName
	}

	if payload.LastName != nil {
		user.LastName = payload.LastName
	} else {
		user.LastName = &userExist.LastName
	}

	if payload.Age != nil {
		user.Age = payload.Age
	} else {
		user.Age = &userExist.Age
	}

	if payload.Email != nil {
		user.Email = payload.Email
	} else {
		user.Email = &userExist.Email
	}

	if payload.RoleId != nil {
		user.RoleId = payload.RoleId
	} else {
		user.RoleId = userExist.RoleId
	}

	if payload.Password != nil {
		hashedPassword, err := utils.HashPassword(*payload.Password)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("problem hashing password"))
			return
		}
		user.Password = &hashedPassword
	} else {
		user.Password = &userExist.Password
	}

	// Handle project associations
	var projectIDs []int
	if payload.ProjectIDS != nil {
		projectIDs = *payload.ProjectIDS
	} else {
		for _, proj := range userExist.Projects {
			projectIDs = append(projectIDs, proj.ID)
		}
	}

	// Update user in the store
	err = h.store.UpdateUser(userId, user, projectIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update user: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	err = h.store.DeleteUser(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleRestoreUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	err = h.store.RestoreUser(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user Id"))
		return
	}
	userId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user Id"))
		return
	}
	now := time.Now()
	err = h.store.UpdateLastActiveTime(userId, now)
}
