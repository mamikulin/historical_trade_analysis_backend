package middleware

import (
	"net/http"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")

        if origin != "" {
            if strings.HasPrefix(origin, "http://localhost") ||
                strings.HasPrefix(origin, "https://localhost") ||
                strings.HasPrefix(origin, "http://127.0.0.1") ||
                strings.HasPrefix(origin, "https://127.0.0.1") ||
                strings.HasPrefix(origin, "http://192.168.") ||
                strings.HasPrefix(origin, "https://192.168.") ||
                strings.HasPrefix(origin, "http://10.") ||
                strings.HasPrefix(origin, "https://10.") ||
                strings.HasPrefix(origin, "http://172.16.") || 
                strings.HasPrefix(origin, "https://172.16.") ||
                strings.HasPrefix(origin, "http://172.20.") || 
                strings.HasPrefix(origin, "tauri://") ||
                strings.HasPrefix(origin, "http://tauri.") ||
                strings.HasPrefix(origin, "https://tauri.") ||
                strings.HasPrefix(origin, "https://mamikulin.github.io"){
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Vary", "Origin")
            }
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, ngrok-skip-browser-warning, Accept, Origin, Referer, User-Agent, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform")
        w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Max-Age", "86400")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}
