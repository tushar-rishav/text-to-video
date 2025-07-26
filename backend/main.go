package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	db     *sql.DB
	rdb    *redis.Client
	logger *logrus.Logger
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Job represents a video generation job
type Job struct {
	ID        string    `json:"job_id"`
	Prompt    string    `json:"prompt"`
	Status    string    `json:"job_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	VideoURL  string    `json:"video_url,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// SubmitRequest represents the request body for job submission
type SubmitRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// SubmitResponse represents the response for job submission
type SubmitResponse struct {
	JobID string `json:"job_id"`
	Error string `json:"error,omitempty"`
}

// StatusResponse represents the response for job status
type StatusResponse struct {
	JobStatus string `json:"job_status"`
	Error     string `json:"error,omitempty"`
}

// ListResponse represents the response for job listing
type ListResponse struct {
	Jobs []Job  `json:"jobs"`
	Error string `json:"error,omitempty"`
}

// VideoResponse represents the response for video retrieval
type VideoResponse struct {
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
}

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Initialize database
	initDatabase()

	// Initialize Redis
	initRedis()
}

func initDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database:", err)
	}

	// Create jobs table if not exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS jobs (
		id VARCHAR(36) PRIMARY KEY,
		prompt TEXT NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		video_url VARCHAR(500),
		error TEXT
	)`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		logger.Fatal("Failed to create jobs table:", err)
	}

	logger.Info("Database initialized successfully")
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis:", err)
	}

	logger.Info("Redis initialized successfully")
}

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/submit", submitJob)
		api.GET("/status", getJobStatus)
		api.GET("/list", listJobs)
		api.GET("/video", getVideo)
		api.GET("/ws", handleWebSocket)
	}

	// Serve static files for videos
	r.Static("/videos", "./videos")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Infof("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}

func submitJob(c *gin.Context) {
	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SubmitResponse{Error: "Invalid request body"})
		return
	}

	jobID := uuid.New().String()
	job := Job{
		ID:        jobID,
		Prompt:    req.Prompt,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert job into database
	_, err := db.Exec("INSERT INTO jobs (id, prompt, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		job.ID, job.Prompt, job.Status, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		logger.Errorf("Failed to insert job: %v", err)
		c.JSON(http.StatusInternalServerError, SubmitResponse{Error: "Failed to create job"})
		return
	}

	// Publish job to Redis for video service
	jobData, _ := json.Marshal(job)
	ctx := context.Background()
	err = rdb.Publish(ctx, "video_jobs", jobData).Err()
	if err != nil {
		logger.Errorf("Failed to publish job to Redis: %v", err)
		// Don't fail the request, just log the error
	}

	logger.Infof("Job submitted successfully: %s", jobID)
	c.JSON(http.StatusOK, SubmitResponse{JobID: jobID})
}

func getJobStatus(c *gin.Context) {
	jobID := c.Query("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, StatusResponse{Error: "job_id is required"})
		return
	}

	var job Job
	err := db.QueryRow("SELECT id, prompt, status, created_at, updated_at, video_url, error FROM jobs WHERE id = ?", jobID).
		Scan(&job.ID, &job.Prompt, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.VideoURL, &job.Error)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, StatusResponse{Error: "Job not found"})
			return
		}
		logger.Errorf("Failed to query job: %v", err)
		c.JSON(http.StatusInternalServerError, StatusResponse{Error: "Failed to get job status"})
		return
	}

	c.JSON(http.StatusOK, StatusResponse{JobStatus: job.Status})
}

func listJobs(c *gin.Context) {
	status := c.Query("status")
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ListResponse{Error: "Invalid offset"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ListResponse{Error: "Invalid limit"})
		return
	}

	var query string
	var args []interface{}

	if status != "" {
		query = "SELECT id, prompt, status, created_at, updated_at, video_url, error FROM jobs WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
		args = []interface{}{status, limit, offset}
	} else {
		query = "SELECT id, prompt, status, created_at, updated_at, video_url, error FROM jobs ORDER BY created_at DESC LIMIT ? OFFSET ?"
		args = []interface{}{limit, offset}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Errorf("Failed to query jobs: %v", err)
		c.JSON(http.StatusInternalServerError, ListResponse{Error: "Failed to list jobs"})
		return
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.Prompt, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.VideoURL, &job.Error)
		if err != nil {
			logger.Errorf("Failed to scan job: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	c.JSON(http.StatusOK, ListResponse{Jobs: jobs})
}

func getVideo(c *gin.Context) {
	jobID := c.Query("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, VideoResponse{Error: "job_id is required"})
		return
	}

	var videoURL string
	err := db.QueryRow("SELECT video_url FROM jobs WHERE id = ? AND status = 'completed'", jobID).Scan(&videoURL)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, VideoResponse{Error: "Video not found or job not completed"})
			return
		}
		logger.Errorf("Failed to query video: %v", err)
		c.JSON(http.StatusInternalServerError, VideoResponse{Error: "Failed to get video"})
		return
	}

	if videoURL == "" {
		c.JSON(http.StatusNotFound, VideoResponse{Error: "Video URL not available"})
		return
	}

	c.JSON(http.StatusOK, VideoResponse{URL: videoURL})
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Subscribe to Redis channel for real-time updates
	pubsub := rdb.Subscribe(context.Background(), "job_updates")
	defer pubsub.Close()

	// Handle incoming messages
	for {
		select {
		case msg := <-pubsub.Channel():
			var job Job
			if err := json.Unmarshal([]byte(msg.Payload), &job); err != nil {
				logger.Errorf("Failed to unmarshal job update: %v", err)
				continue
			}

			if err := conn.WriteJSON(job); err != nil {
				logger.Errorf("Failed to send job update: %v", err)
				return
			}
		}
	}
} 