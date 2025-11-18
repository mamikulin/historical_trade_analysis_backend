package user

import (
	"archpath/internal/middleware"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	service    *Service
	jwtSecret  string
	jwtExpiry  time.Duration
}

func NewHandler(service *Service, jwtSecret string, jwtExpiry time.Duration) *Handler {
	return &Handler{
		service:    service,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
	}
}

func (h *Handler) RegisterRoutes(r interface{}) {
	// Этот метод будет вызван из api/server.go
}

// @Summary Register a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body object{user=User,password=string} true "Registration data"
// @Success 201 {object} User
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Registration failed"
// @Router /users/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		User     *User  `json:"user"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.RegisterUser(requestData.User, requestData.Password)
	if err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	requestData.User.PasswordHash = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(requestData.User)
}

// @Summary Login
// @Description Authenticate user and get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body object{login=string,password=string} true "Login credentials"
// @Success 200 {object} object{user=User,token=string}
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Authentication failed"
// @Router /users/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid login request", http.StatusBadRequest)
		return
	}

	user, err := h.service.AuthenticateUser(loginData.Login, loginData.Password)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Генерируем JWT токен
	token, err := middleware.GenerateJWT(user.ID, user.Role, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""
	
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Security BearerAuth
// @Router /api/users/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.PasswordHash = ""
	json.NewEncoder(w).Encode(user)
}

// @Summary Update current user
// @Description Update information of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 200 {object} User
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Update failed"
// @Security BearerAuth
// @Router /users/me [put]
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Запрещаем изменение критичных полей
	user.ID = 0
	user.PasswordHash = ""
	user.Role = ""

	err = h.service.UpdateUser(userID, &user)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	updatedUser, _ := h.service.GetUserByID(userID)
	updatedUser.PasswordHash = ""
	json.NewEncoder(w).Encode(updatedUser)
}

// @Summary Logout
// @Description Logout current user (client-side token deletion)
// @Tags users
// @Produce json
// @Success 200 {object} object{message=string}
// @Security BearerAuth
// @Router /users/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// With JWT, logout is handled client-side by deleting the token
	// Optionally, you could implement a token blacklist here
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully. Please delete your token.",
	})
}