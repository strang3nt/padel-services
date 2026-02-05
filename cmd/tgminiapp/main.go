package main

import (
	"bytes"
	"context"
	"embed"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/strang3nt/padel-services/pkg/services"
	"github.com/strang3nt/padel-services/pkg/tournament"

	"github.com/strang3nt/padel-services/pkg/database"

	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ctx                    = context.Background()
	botToken               = os.Getenv("BOT_TOKEN")
	jwtSecret              = []byte(os.Getenv("JWT_SECRET"))
	database_url           = os.Getenv("DATABASE_URL")
	downloadTokens         = makeFileTokenHandler()
	tournamentPdfGenerator = services.MakeTournamentPdfGenerator()
)

type AuthRequest struct {
	InitData string `json:"initDataRaw" binding:"required"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["sub"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
	}
}

//go:embed dist/*
var frontendFiles embed.FS

type Response struct {
	Message string `json:"message"`
}

func main() {
	r := gin.Default()

	conn, err := pgxpool.New(ctx, database_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	users, err := database.GetUsersIds(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to retrieve allowed users: %v\n", err)
		os.Exit(1)
	}
	whitelistedIDs := MakeAllowerUsersFromUserData(users)

	r.POST("/auth", func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Print(err)
			c.JSON(400, gin.H{"error": "Payload missing"})
			return
		}

		telegramID, err := verifyTelegram(req.InitData)
		if err != nil {
			c.JSON(401, gin.H{"error": "Tampered data"})
			return
		}

		// Whitelist check
		if !whitelistedIDs.IsUserAllowed(telegramID) {
			c.JSON(403, gin.H{"error": "Access denied: ID not whitelisted"})
			return
		}

		// Issue JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": telegramID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		})
		tokenString, _ := token.SignedString(jwtSecret)

		c.JSON(200, gin.H{"token": tokenString, "id": telegramID})
	})
	// Protected Routes
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{

		protected.GET("/tournaments", func(c *gin.Context) {
			// DefaultQuery provides a fallback if the key doesn't exist
			dateString := c.Query("date")
			date, _ := time.Parse("2006-01-02", dateString)
			userIdBlob, _ := c.Get("user_id")
			userId, _ := userIdBlob.(int64)

			tournaments, _ := database.GetTournamentsByDate(ctx, conn, userId, date)
			c.JSON(http.StatusOK, gin.H{
				"date":        dateString,
				"tournaments": tournaments})
		})

		protected.POST("/create-tournament", func(c *gin.Context) {
			fmt.Print("Create tournament handler called")
			tournamentType := c.Query("tournamentType")
			dateStart, _ := time.Parse(time.RFC3339, c.Query("dateStart"))
			totalRounds, _ := strconv.ParseInt(c.Query("totalRounds"), 10, 32)
			availableCourts, _ := strconv.ParseInt(c.Query("availableCourts"), 10, 32)
			userIdBlob, _ := c.Get("user_id")
			userId, _ := userIdBlob.(int64)
			var teams []tournament.Team
			if err := c.ShouldBindJSON(&teams); err != nil {
				fmt.Print(err)
				c.JSON(400, gin.H{"error": "Payload missing"})
				return
			}

			tournament := services.CreateTournament(tournamentType, dateStart, teams, int(totalRounds), int(availableCourts))

			if tournament != nil {
				err = database.CreateTournament(ctx, conn, userId, &tournament)
				if err != nil {
					c.JSON(500, gin.H{"error": "could not save tournament"})
					return
				}
				c.Status(200)
			} else {
				c.JSON(400, gin.H{"error": "could not create tournament"})
			}
		})

		protected.POST("/tournament/generate-link", func(c *gin.Context) {
			var req tournament.TournamentData
			if err := c.ShouldBindJSON(&req); err != nil {
				fmt.Print(err)
				c.JSON(400, gin.H{"error": "Payload missing"})
				return
			}
			token := downloadTokens.GenerateToken(5*time.Minute, req)
			c.JSON(200, gin.H{
				"token": token,
			})
		})

	}

	r.GET("/api/tournament/download", func(c *gin.Context) {
		token := c.Query("token")
		c.Header("Access-Control-Allow-Origin", "https://web.telegram.org")

		blob, exists := downloadTokens.GetBlob(token)
		tournament := blob.(tournament.TournamentData)
		if !exists {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		pdfPath, err := tournamentPdfGenerator.CreatePdfTournament(
			services.FromTournamentDataToTemplateData(tournament),
			services.Rodeo,
			fmt.Sprint(tournament.Date.Format("2006-01-02"), "_", tournament.Name, "_", token),
		)
		if err != nil {
			log.Printf("Error creating PDF: %v", err)
			return
		}

		defer func() {
			err := os.Remove(pdfPath)
			if err != nil {
				log.Printf("error while removing file: %v", err)
			}
		}()

		c.FileAttachment(
			pdfPath,
			filepath.Base(pdfPath),
		)
	})

	publicFiles, _ := fs.Sub(frontendFiles, "dist")
	staticServer := http.FS(publicFiles)

	// 3. THE SPA SERVING LOGIC
	r.NoRoute(func(c *gin.Context) {

		path := c.Request.URL.Path

		// If the request is for a real file (like /assets/index-123.js), serve it
		_, err := publicFiles.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			http.FileServer(staticServer).ServeHTTP(c.Writer, c.Request)
			return
		}

		// Otherwise, it's a React Router path (like /dashboard).
		// Serve index.html and let React handle the rest.
		index, _ := publicFiles.Open("index.html")
		stat, _ := index.Stat()
		content, _ := io.ReadAll(index)
		http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), bytes.NewReader(content))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
