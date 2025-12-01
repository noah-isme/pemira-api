package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"pemira-api/internal/http/response"
	"pemira-api/internal/shared/ctxkeys"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid.")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		response.UnprocessableEntity(w, "VALIDATION_ERROR", "username dan password wajib diisi.")
		return
	}

	// Extract user agent and IP
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.Header.Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
		if ipAddress != "" {
			// Take first IP if multiple
			ipAddress = strings.Split(ipAddress, ",")[0]
		}
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// Remove port from IP address if present
	if idx := strings.LastIndex(ipAddress, ":"); idx != -1 {
		// Check if it's IPv6 or IPv4 with port
		if strings.Count(ipAddress, ":") == 1 || strings.HasPrefix(ipAddress, "[") {
			ipAddress = strings.TrimRight(strings.Split(ipAddress, ":")[0], "]")
			ipAddress = strings.TrimLeft(ipAddress, "[")
		}
	}

	loginResp, err := h.service.Login(r.Context(), req, userAgent, ipAddress)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, loginResp)
}

// RegisterStudent handles POST /auth/register/student
func (h *AuthHandler) RegisterStudent(w http.ResponseWriter, r *http.Request) {
	var req RegisterStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid.")
		return
	}

	user, err := h.service.RegisterStudent(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":        user,
		"message":     "Registrasi mahasiswa berhasil.",
		"voting_mode": normalizeVotingMode(req.VotingMode),
	})
}

// RegisterLecturerStaff handles POST /auth/register/lecturer-staff
func (h *AuthHandler) RegisterLecturerStaff(w http.ResponseWriter, r *http.Request) {
	var req RegisterLecturerStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid.")
		return
	}

	user, err := h.service.RegisterLecturerStaff(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":        user,
		"message":     "Registrasi berhasil.",
		"voting_mode": normalizeVotingMode(req.VotingMode),
	})
}

// RefreshToken handles POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid.")
		return
	}

	if req.RefreshToken == "" {
		response.UnprocessableEntity(w, "VALIDATION_ERROR", "Refresh token wajib diisi.")
		return
	}

	refreshResp, err := h.service.RefreshToken(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, refreshResp)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid.")
		return
	}

	if req.RefreshToken == "" {
		response.UnprocessableEntity(w, "VALIDATION_ERROR", "Refresh token wajib diisi.")
		return
	}

	if err := h.service.Logout(r.Context(), req); err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully.",
	})
}

// Me handles GET /auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := ctxkeys.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "UNAUTHORIZED", "Token tidak valid atau tidak ditemukan.")
		return
	}

	authUser, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, authUser)
}

// handleError maps service errors to HTTP responses
func (h *AuthHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		response.Unauthorized(w, "INVALID_CREDENTIALS", "Username atau password salah.")

	case errors.Is(err, ErrInactiveUser):
		response.Forbidden(w, "USER_INACTIVE", "Akun tidak aktif.")

	case errors.Is(err, ErrInvalidRefreshToken):
		response.Unauthorized(w, "INVALID_REFRESH_TOKEN", "Refresh token tidak valid atau sudah kadaluarsa.")

	case errors.Is(err, ErrUserNotFound):
		response.NotFound(w, "USER_NOT_FOUND", "Pengguna tidak ditemukan.")

	case errors.Is(err, ErrUsernameExists):
		response.Conflict(w, "USERNAME_EXISTS", "Username sudah terdaftar.")

	case errors.Is(err, ErrNIMExists):
		response.Conflict(w, "NIM_EXISTS", "NIM sudah terdaftar.")

	case errors.Is(err, ErrNIDNExists):
		response.Conflict(w, "NIDN_EXISTS", "NIDN sudah terdaftar.")

	case errors.Is(err, ErrNIPExists):
		response.Conflict(w, "NIP_EXISTS", "NIP sudah terdaftar.")

	case errors.Is(err, ErrInvalidRegisterType):
		response.BadRequest(w, "INVALID_REQUEST", "Tipe registrasi tidak valid.")

	case errors.Is(err, ErrInvalidRegistration):
		response.UnprocessableEntity(w, "VALIDATION_ERROR", "Data registrasi tidak lengkap atau tidak valid.")

	case errors.Is(err, ErrModeNotAvailable):
		response.UnprocessableEntity(w, "MODE_NOT_AVAILABLE", "Mode tidak tersedia untuk pemilu ini.")

	default:
		// Log internal error
		slog.Error("auth handler error", "error", err)
		response.InternalServerError(w, "INTERNAL_ERROR", "Terjadi kesalahan pada sistem.")
	}
}

// LogoutPage handles GET /auth/logout-page - simple HTML page to clear tokens
func (h *AuthHandler) LogoutPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Logout - PEMIRA</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            text-align: center;
            max-width: 400px;
        }
        h1 { color: #333; margin-bottom: 1rem; }
        p { color: #666; margin-bottom: 1.5rem; }
        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #667eea;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto 1rem;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        .success { color: #10b981; font-weight: bold; }
        a {
            display: inline-block;
            margin-top: 1rem;
            padding: 0.5rem 2rem;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 5px;
        }
        a:hover { background: #5568d3; }
    </style>
</head>
<body>
    <div class="container">
        <div class="spinner" id="spinner"></div>
        <h1>Logout</h1>
        <p id="message">Menghapus sesi Anda...</p>
        <div id="loginLink" style="display:none;">
            <p class="success">Berhasil logout!</p>
            <a href="/">Kembali ke Login</a>
        </div>
    </div>
    <script>
        // Clear all possible token storage locations
        localStorage.clear();
        sessionStorage.clear();
        
        // Clear specific token keys (in case they use specific names)
        const tokenKeys = ['token', 'access_token', 'refresh_token', 'auth_token', 'jwt'];
        tokenKeys.forEach(key => {
            localStorage.removeItem(key);
            sessionStorage.removeItem(key);
        });
        
        // Wait a moment then show success
        setTimeout(() => {
            document.getElementById('spinner').style.display = 'none';
            document.getElementById('message').style.display = 'none';
            document.getElementById('loginLink').style.display = 'block';
        }, 1000);
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
