package store

import (
	"database/sql"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storer interface {
	CreateUser(username, passwordHash string) (int64, error)
	GetUserByUsername(username string) (*User, error)
	CreateHabit(userID int64, name string) (int64, error)
	GetHabitsByUserID(userID int64) ([]Habit, error)
	CreateHabitTrack(userID, habitID int64) (int64, error)
	GetHabitTrackForToday(userID, habitID int64) (*HabitTrack, error)
	DeleteHabitTrack(trackID int64) error
	GetHabitStreak(userID, habitID int64) (int, error)
	GetHabitCompletionRate(userID, habitID int64, windowDays int) (float64, error)
	GetAverageSleepDuration(userID int64, windowDays int) (float64, error)
	GetAverageSleepQuality(userID int64, windowDays int) (float64, error)
	GetAverageDailyCalories(userID int64, windowDays int) (float64, error)
	GetNutritionLeaderboard(windowDays int) ([]LeaderboardEntry, error)
	GetSleepLeaderboard(windowDays int) ([]LeaderboardEntry, error)
	GetHabitLeaderboard(windowDays int) ([]LeaderboardEntry, error)
	CreateNutritionLog(userID int64, mealName string, calories int) (int64, error)
	GetNutritionLogsByUserID(userID int64) ([]NutritionLog, error)
	CreateSleepLog(userID int64, duration float64, quality int) (int64, error)
	GetSleepLogsByUserID(userID int64) ([]SleepLog, error)
	CreateMoodLog(userID int64, rating int, notes string) (int64, error)
	UpdateMoodLogAnalysis(logID int64, analysis string) error
	GetRecentNutritionLogs(userID int64, windowDays int) ([]NutritionLog, error)
	GetMoodLogsByUserID(userID int64) ([]MoodLog, error)
}

type Store struct {
	db *sql.DB
}

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

type Habit struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
}

type HabitTrack struct {
	ID        int64
	HabitID   int64
	UserID    int64
	TrackedAt time.Time
}

type NutritionLog struct {
	ID        int64
	UserID    int64
	MealName  string
	Calories  int
	CreatedAt time.Time
}

type SleepLog struct {
	ID            int64
	UserID        int64
	DurationHours float64
	QualityRating int
	CreatedAt     time.Time
}

type MoodLog struct {
	ID         int64
	UserID     int64
	MoodRating int
	Notes      string
	AIAnalysis string `json:"ai_analysis,omitempty"`
	CreatedAt  time.Time
}

type LeaderboardEntry struct {
	Username string  `json:"username"`
	Score    float64 `json:"score"`
}

func NewStore() (Storer, error) {
	db, err := sql.Open("sqlite3", "./habits.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func createTables(db *sql.DB) error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);
	`
	habitsTable := `
	CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	habitTracksTable := `
	CREATE TABLE IF NOT EXISTS habit_tracks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		habit_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		tracked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (habit_id) REFERENCES habits(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	nutritionTable := `
	CREATE TABLE IF NOT EXISTS nutrition_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		meal_name TEXT NOT NULL,
		calories INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	sleepTable := `
	CREATE TABLE IF NOT EXISTS sleep_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		duration_hours REAL NOT NULL,
		quality_rating INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	moodTable := `
	CREATE TABLE IF NOT EXISTS mood_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		mood_rating INTEGER NOT NULL,
		notes TEXT,
		ai_analysis TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	for _, table := range []string{usersTable, habitsTable, habitTracksTable, nutritionTable, sleepTable, moodTable} {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) CreateUser(username, passwordHash string) (int64, error) {
	res, err := s.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	row := s.db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) CreateHabit(userID int64, name string) (int64, error) {
	res, err := s.db.Exec("INSERT INTO habits (user_id, name) VALUES (?, ?)", userID, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetHabitsByUserID(userID int64) ([]Habit, error) {
	rows, err := s.db.Query("SELECT id, user_id, name, created_at FROM habits WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit
		if err := rows.Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.CreatedAt); err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}

	return habits, nil
}

func (s *Store) CreateHabitTrack(userID, habitID int64) (int64, error) {
	res, err := s.db.Exec("INSERT INTO habit_tracks (user_id, habit_id) VALUES (?, ?)", userID, habitID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetHabitTrackForToday(userID, habitID int64) (*HabitTrack, error) {
	track := &HabitTrack{}
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	row := s.db.QueryRow("SELECT id, habit_id, user_id, tracked_at FROM habit_tracks WHERE user_id = ? AND habit_id = ? AND tracked_at >= ? AND tracked_at < ?", userID, habitID, today, tomorrow)
	if err := row.Scan(&track.ID, &track.HabitID, &track.UserID, &track.TrackedAt); err != nil {
		return nil, err
	}
	return track, nil
}

func (s *Store) DeleteHabitTrack(trackID int64) error {
	_, err := s.db.Exec("DELETE FROM habit_tracks WHERE id = ?", trackID)
	return err
}

func (s *Store) GetHabitStreak(userID, habitID int64) (int, error) {
	rows, err := s.db.Query("SELECT tracked_at FROM habit_tracks WHERE user_id = ? AND habit_id = ? ORDER BY tracked_at DESC", userID, habitID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return 0, err
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0, nil
	}

	uniqueDates := make(map[time.Time]bool)
	var distinctDates []time.Time
	for _, date := range dates {
		day := date.Truncate(24 * time.Hour)
		if !uniqueDates[day] {
			uniqueDates[day] = true
			distinctDates = append(distinctDates, day)
		}
	}
	sort.Slice(distinctDates, func(i, j int) bool {
		return distinctDates[i].After(distinctDates[j])
	})

	streak := 0
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	if len(distinctDates) > 0 && (distinctDates[0] == today || distinctDates[0] == yesterday) {
		streak = 1
		for i := 0; i < len(distinctDates)-1; i++ {
			day1 := distinctDates[i]
			day2 := distinctDates[i+1]
			if day1.Sub(day2).Hours() == 24 {
				streak++
			} else {
				break
			}
		}
	}

	return streak, nil
}

func (s *Store) GetHabitCompletionRate(userID, habitID int64, windowDays int) (float64, error) {
	var count int
	startTime := time.Now().AddDate(0, 0, -windowDays)

	row := s.db.QueryRow("SELECT COUNT(DISTINCT date(tracked_at)) FROM habit_tracks WHERE user_id = ? AND habit_id = ? AND tracked_at >= ?", userID, habitID, startTime)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return float64(count) / float64(windowDays) * 100, nil
}

func (s *Store) GetAverageSleepDuration(userID int64, windowDays int) (float64, error) {
	var avg sql.NullFloat64
	startTime := time.Now().AddDate(0, 0, -windowDays)

	row := s.db.QueryRow("SELECT AVG(duration_hours) FROM sleep_logs WHERE user_id = ? AND created_at >= ?", userID, startTime)
	if err := row.Scan(&avg); err != nil {
		return 0, err
	}
	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}

func (s *Store) GetAverageSleepQuality(userID int64, windowDays int) (float64, error) {
	var avg sql.NullFloat64
	startTime := time.Now().AddDate(0, 0, -windowDays)

	row := s.db.QueryRow("SELECT AVG(quality_rating) FROM sleep_logs WHERE user_id = ? AND created_at >= ?", userID, startTime)
	if err := row.Scan(&avg); err != nil {
		return 0, err
	}
	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}

func (s *Store) GetAverageDailyCalories(userID int64, windowDays int) (float64, error) {
	var totalCalories sql.NullFloat64
	startTime := time.Now().AddDate(0, 0, -windowDays)

	row := s.db.QueryRow("SELECT SUM(calories) FROM nutrition_logs WHERE user_id = ? AND created_at >= ?", userID, startTime)
	if err := row.Scan(&totalCalories); err != nil {
		return 0, err
	}
	if !totalCalories.Valid {
		return 0, nil
	}
	return totalCalories.Float64 / float64(windowDays), nil
}

func (s *Store) GetNutritionLeaderboard(windowDays int) ([]LeaderboardEntry, error) {
	startTime := time.Now().AddDate(0, 0, -windowDays)
	query := `
		SELECT u.username, COUNT(n.id) as score
		FROM nutrition_logs n
		JOIN users u ON n.user_id = u.id
		WHERE n.created_at >= ?
		GROUP BY u.username
		ORDER BY score DESC
		LIMIT 10
	`
	rows, err := s.db.Query(query, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.Username, &entry.Score); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}
	return leaderboard, nil
}

func (s *Store) GetSleepLeaderboard(windowDays int) ([]LeaderboardEntry, error) {
	startTime := time.Now().AddDate(0, 0, -windowDays)
	query := `
		SELECT u.username, SUM(s.duration_hours) as score
		FROM sleep_logs s
		JOIN users u ON s.user_id = u.id
		WHERE s.created_at >= ?
		GROUP BY u.username
		ORDER BY score DESC
		LIMIT 10
	`
	rows, err := s.db.Query(query, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.Username, &entry.Score); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}
	return leaderboard, nil
}

func (s *Store) GetHabitLeaderboard(windowDays int) ([]LeaderboardEntry, error) {
	startTime := time.Now().AddDate(0, 0, -windowDays)
	query := `
		SELECT u.username, COUNT(ht.id) as score
		FROM habit_tracks ht
		JOIN users u ON ht.user_id = u.id
		WHERE ht.tracked_at >= ?
		GROUP BY u.username
		ORDER BY score DESC
		LIMIT 10
	`
	rows, err := s.db.Query(query, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.Username, &entry.Score); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}
	return leaderboard, nil
}

func (s *Store) CreateNutritionLog(userID int64, mealName string, calories int) (int64, error) {
	res, err := s.db.Exec("INSERT INTO nutrition_logs (user_id, meal_name, calories) VALUES (?, ?, ?)", userID, mealName, calories)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetNutritionLogsByUserID(userID int64) ([]NutritionLog, error) {
	rows, err := s.db.Query("SELECT id, user_id, meal_name, calories, created_at FROM nutrition_logs WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []NutritionLog
	for rows.Next() {
		var log NutritionLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.MealName, &log.Calories, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *Store) CreateSleepLog(userID int64, duration float64, quality int) (int64, error) {
	res, err := s.db.Exec("INSERT INTO sleep_logs (user_id, duration_hours, quality_rating) VALUES (?, ?, ?)", userID, duration, quality)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetSleepLogsByUserID(userID int64) ([]SleepLog, error) {
	rows, err := s.db.Query("SELECT id, user_id, duration_hours, quality_rating, created_at FROM sleep_logs WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []SleepLog
	for rows.Next() {
		var log SleepLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.DurationHours, &log.QualityRating, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *Store) CreateMoodLog(userID int64, rating int, notes string) (int64, error) {
	res, err := s.db.Exec("INSERT INTO mood_logs (user_id, mood_rating, notes) VALUES (?, ?, ?)", userID, rating, notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateMoodLogAnalysis(logID int64, analysis string) error {
	_, err := s.db.Exec("UPDATE mood_logs SET ai_analysis = ? WHERE id = ?", analysis, logID)
	return err
}

func (s *Store) GetRecentNutritionLogs(userID int64, windowDays int) ([]NutritionLog, error) {
	startTime := time.Now().AddDate(0, 0, -windowDays)
	rows, err := s.db.Query("SELECT id, user_id, meal_name, calories, created_at FROM nutrition_logs WHERE user_id = ? AND created_at >= ? ORDER BY created_at DESC", userID, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []NutritionLog
	for rows.Next() {
		var log NutritionLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.MealName, &log.Calories, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *Store) GetMoodLogsByUserID(userID int64) ([]MoodLog, error) {
	rows, err := s.db.Query("SELECT id, user_id, mood_rating, notes, ai_analysis, created_at FROM mood_logs WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []MoodLog
	for rows.Next() {
		var log MoodLog
		var analysis sql.NullString
		if err := rows.Scan(&log.ID, &log.UserID, &log.MoodRating, &log.Notes, &analysis, &log.CreatedAt); err != nil {
			return nil, err
		}
		if analysis.Valid {
			log.AIAnalysis = analysis.String
		}
		logs = append(logs, log)
	}

	return logs, nil
}
