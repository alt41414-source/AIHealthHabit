package api

import (
	"aiHealthHabit/backend/ai"
	"aiHealthHabit/backend/store"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type API struct {
	store store.Storer
	ai    ai.AIClient
}

func New(store store.Storer, ai ai.AIClient) *API {
	return &API{store: store, ai: ai}
}

func (a *API) TrackHabit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)
	habitID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	track, err := a.store.GetHabitTrackForToday(userID, habitID)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if track != nil {
		if err := a.store.DeleteHabitTrack(track.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "untracked"})
	} else {
		_, err := a.store.CreateHabitTrack(userID, habitID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "tracked"})
	}
}

func (a *API) GetHabitStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)
	habitID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	streak, err := a.store.GetHabitStreak(userID, habitID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	completionRate, err := a.store.GetHabitCompletionRate(userID, habitID, 30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"streak":          streak,
		"completion_rate": completionRate,
	})
}

func (a *API) GetSleepStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	avgDuration, err := a.store.GetAverageSleepDuration(userID, 30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	avgQuality, err := a.store.GetAverageSleepQuality(userID, 30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"average_duration": avgDuration,
		"average_quality":  avgQuality,
	})
}

func (a *API) GetNutritionStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	avgCalories, err := a.store.GetAverageDailyCalories(userID, 30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"average_daily_calories": avgCalories,
	})
}

func (a *API) GetNutritionLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := a.store.GetNutritionLeaderboard(30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(leaderboard)
}

func (a *API) GetSleepLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := a.store.GetSleepLeaderboard(30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(leaderboard)
}

func (a *API) GetHabitLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := a.store.GetHabitLeaderboard(30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(leaderboard)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = a.store.CreateUser(req.Username, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.store.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{Token: tokenString, Username: user.Username}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type CreateHabitRequest struct {
	Name string `json:"name"`
}

func (a *API) CreateHabit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	var req CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := a.store.CreateHabit(userID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) GetHabits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	habits, err := a.store.GetHabitsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

type CreateNutritionLogRequest struct {
	MealName string `json:"meal_name"`
	Calories int    `json:"calories"`
}

func (a *API) CreateNutritionLog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	var req CreateNutritionLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := a.store.CreateNutritionLog(userID, req.MealName, req.Calories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) GetNutritionLogs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	logs, err := a.store.GetNutritionLogsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

type CreateSleepLogRequest struct {
	DurationHours float64 `json:"duration_hours"`
	QualityRating int     `json:"quality_rating"`
}

func (a *API) CreateSleepLog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	var req CreateSleepLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := a.store.CreateSleepLog(userID, req.DurationHours, req.QualityRating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) GetSleepLogs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	logs, err := a.store.GetSleepLogsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

type CreateMoodLogRequest struct {
	MoodRating int    `json:"mood_rating"`
	Notes      string `json:"notes"`
}

func (a *API) CreateMoodLog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	var req CreateMoodLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logID, err := a.store.CreateMoodLog(userID, req.MoodRating, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		if req.Notes != "" {
			analysis, err := a.ai.AnalyzeText(context.Background(), req.Notes, []string{"openai/gpt-oss-120b", "openai/gpt-oss-20b", "moonshotai/kimi-k2-instruct-0905", "llama-3.3-70b-versatile", "meta-llama/llama-4-maverick-17b-128e-instruct", "meta-llama/llama-4-scout-17b-16e-instruct", "groq/compound", "groq/compound-mini", "qwen/qwen3-32b"})
			if err != nil {
				log.Printf("Error analyzing mood log: %v", err)
				return
			}
			if err := a.store.UpdateMoodLogAnalysis(logID, analysis); err != nil {
				log.Printf("Error updating mood log analysis: %v", err)
			}
		}
	}()

	w.WriteHeader(http.StatusCreated)
}

func (a *API) GetMoodLogs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	logs, err := a.store.GetMoodLogsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (a *API) GetNutritionSuggestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	logs, err := a.store.GetRecentNutritionLogs(userID, 7)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(logs) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"suggestion": "Not enough data to provide a suggestion. Please log your meals for a few days."})
		return
	}

	var promptBuilder strings.Builder
	promptBuilder.WriteString("Here are my recent meals:\n")
	for _, log := range logs {
		promptBuilder.WriteString("- ")
		promptBuilder.WriteString(log.MealName)
		promptBuilder.WriteString(" (")
		promptBuilder.WriteString(strconv.Itoa(log.Calories))
		promptBuilder.WriteString(" calories)\n")
	}
	promptBuilder.WriteString("\nBased on these meals, what is one simple, healthy, and actionable suggestion you can give me for my next meal? Keep it concise and positive.")

	suggestion, err := a.ai.AnalyzeText(r.Context(), promptBuilder.String(), []string{"openai/gpt-oss-120b", "openai/gpt-oss-20b", "moonshotai/kimi-k2-instruct-0905", "llama-3.3-70b-versatile", "meta-llama/llama-4-maverick-17b-128e-instruct", "meta-llama/llama-4-scout-17b-16e-instruct", "groq/compound", "groq/compound-mini", "qwen/qwen3-32b"})
	if err != nil {
		http.Error(w, "Failed to get suggestion from AI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"suggestion": suggestion})
}
