package user

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

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
    
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    err = h.service.RegisterUser(requestData.User, requestData.Password)
    if err != nil {
        http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    requestData.User.PasswordHash = ""
    json.NewEncoder(w).Encode(requestData.User)
}

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

    user.PasswordHash = ""
    json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
    user := &User{
        Login: "user1",  

    json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err = h.service.UpdateUser(1, &user)
    if err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
