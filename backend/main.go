package main

import (
	"aiHealthHabit/backend/ai"
	"aiHealthHabit/backend/api"
	"aiHealthHabit/backend/store"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	storage, err := store.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	aiClient := ai.NewClient()
	apiHandler := api.New(storage, aiClient)

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", apiHandler.Register)
		r.Post("/login", apiHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(api.AuthMiddleware)

			r.Post("/habits", apiHandler.CreateHabit)
			r.Get("/habits", apiHandler.GetHabits)
			r.Post("/habits/{id}/track", apiHandler.TrackHabit)
			r.Get("/habits/{id}/stats", apiHandler.GetHabitStats)

			r.Post("/nutrition", apiHandler.CreateNutritionLog)
			r.Get("/nutrition", apiHandler.GetNutritionLogs)
			r.Get("/nutrition/stats", apiHandler.GetNutritionStats)
			r.Get("/nutrition/suggestion", apiHandler.GetNutritionSuggestion)

			r.Post("/sleep", apiHandler.CreateSleepLog)
			r.Get("/sleep", apiHandler.GetSleepLogs)
			r.Get("/sleep/stats", apiHandler.GetSleepStats)

			r.Post("/mood", apiHandler.CreateMoodLog)
			r.Get("/mood", apiHandler.GetMoodLogs)

			r.Get("/leaderboard/nutrition", apiHandler.GetNutritionLeaderboard)
			r.Get("/leaderboard/sleep", apiHandler.GetSleepLeaderboard)
			r.Get("/leaderboard/habit", apiHandler.GetHabitLeaderboard)
		})
	})

	workDir, _ := os.Getwd()
	frontendDir := filepath.Join(workDir, "..", "frontend", "dist")

	fs := http.FileServer(http.Dir(frontendDir))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(filepath.Join(frontendDir, r.URL.Path)); !os.IsNotExist(err) || r.URL.Path == "/" {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, r)
}
