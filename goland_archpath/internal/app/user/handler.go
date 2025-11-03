package user

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

// Handler for the User model
type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/users/register", h.Register).Methods("POST")
    r.HandleFunc("/users/login", h.Login).Methods("POST")
    r.HandleFunc("/users/me", h.GetMe).Methods("GET")
    r.HandleFunc("/users/me", h.UpdateMe).Methods("PUT")
    r.HandleFunc("/users/logout", h.Logout).Methods("POST")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        User     *User  `json:"user"`
        Password string `json:"password"`
    }
    
    // Decode the incoming JSON request body
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Register the user
    err = h.service.RegisterUser(requestData.User, requestData.Password)
    if err != nil {
        http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Respond with the user data (excluding password hash)
    requestData.User.PasswordHash = ""
    json.NewEncoder(w).Encode(requestData.User)
}

// Authenticate user (POST /users/login)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var loginData struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }

    // Decode login request
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

    // Respond with the user data (excluding password hash)
    user.PasswordHash = ""
    json.NewEncoder(w).Encode(user)
}

// Get user data after authentication (GET /users/me)
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
    // This is a simplified example; typically you'd extract user info from a session or JWT token.
    user := &User{
        Login: "user1",  // Example: hardcoded, should be retrieved from session or token
    }

    json.NewEncoder(w).Encode(user)
}

// Update user info (PUT /users/me)
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Update user details
    err = h.service.UpdateUser(1, &user) // Assuming user ID is 1 for simplicity
    if err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

// Logout (POST /users/logout)
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
    // In a real application, you'd clear sessions or JWT tokens here.
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
