package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type logoData struct {
	LogoId       int64
	ByteCode     []byte
	SportsCenter string
	MimeType     string
	UserIDs      []int64
}

type UserData struct {
	Id           int64
	SportsCentre string
	Logo         string
}

func GetUsersIds(ctx context.Context, conn *pgxpool.Pool) ([]UserData, error) {

	query := `
SELECT
		l.id
    l.logo,
		sc.name,
    l.mime_type, 
    array_agg(u.id) AS user_ids
FROM logos l
JOIN sports_center sc ON sc.logo_id = l.id
JOIN users u ON u.sports_center_id = sc.id
GROUP BY l.id;`

	downloadDir := "static/logos"
	os.MkdirAll(downloadDir, 0755)

	rows, _ := conn.Query(ctx, query)
	usersLogos, _ := pgx.CollectRows(rows, pgx.RowToStructByPos[logoData])
	users := make([]UserData, 0)
	for _, x := range usersLogos {
		ext := ".bin"
		switch x.MimeType {
		case "image/svg+xml":
			ext = ".svg"
		case "image/png":
			ext = ".png"
		}

		fileName := fmt.Sprintf("logo_%d%s", x.LogoId, ext)
		relativePath := filepath.Join(downloadDir, fileName)

		err := os.WriteFile(relativePath, x.ByteCode, 0644)
		if err != nil {
			log.Printf("Failed to save %s: %v", relativePath, err)
			continue
		}
		for _, u := range x.UserIDs {
			users = append(users, UserData{
				Id:           u,
				SportsCentre: x.SportsCenter,
				Logo:         relativePath,
			})
		}
	}

	return users, nil

}
