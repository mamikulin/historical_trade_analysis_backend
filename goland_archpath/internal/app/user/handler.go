package user

import (
	"archpath/internal/app/session"
	"archpath/internal/middleware"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service        *Service
	sessionManager *session.Manager
}

func NewHandler(service *Service, sessionManager *session.Manager) *Handler {
	return &Handler{
		service:        service,
		sessionManager: sessionManager,
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
// @Description Authenticate user and create session
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body object{login=string,password=string} true "Login credentials"
// @Success 200 {object} User
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

	// Создаем сессию
	sessionID, err := h.sessionManager.CreateSession(r.Context(), user.ID, user.Role)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // В продакшене должно быть true
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 часа
	})

	user.PasswordHash = ""
	json.NewEncoder(w).Encode(user)
}

// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Security CookieAuth
// @Router /users/me [get]
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
// @Security CookieAuth
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
// @Description Logout current user and destroy session
// @Tags users
// @Produce json
// @Success 200 {object} object{message=string}
// @Security CookieAuth
// @Router /users/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Удаляем сессию из Redis
		h.sessionManager.DeleteSession(r.Context(), cookie.Value)
	}

	// Очищаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}