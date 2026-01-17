package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"padelservices/pkg/tournament"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func queryInsertTeam(ctx context.Context, tx pgx.Tx, team1 tournament.Team) (int64, error) {

	const sql = `
    WITH upserted_people AS (
        INSERT INTO person (name)
        VALUES ($1), ($2)
        ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
        RETURNING id, name
    )
    INSERT INTO team (person1_id, person2_id)
    SELECT 
        (SELECT id FROM upserted_people WHERE name = $1),
        (SELECT id FROM upserted_people WHERE name = $2)
    RETURNING id;`

	person1 := team1.Person_1.Id
	person2 := team1.Person_2.Id

	var id int64
	if err := tx.QueryRow(ctx, sql, person1, person2).Scan(&id); err != nil {
		return -1, fmt.Errorf("error while inserting team: %w", err)
	}

	return id, nil
}

func queryCreateMatch(ctx context.Context, tx pgx.Tx, round_number int, tournamentId, team1Id, team2Id int64, court_number int) error {

	const sql = `
		WITH

		new_match AS (
				INSERT INTO match (team1_id, team2_id, court_number)
				VALUES (
						($1),	($2), ($3)
				)
				RETURNING id
		)

		INSERT INTO round_tournament (tournament_id, match_id, round_number)
		VALUES (
				($4),
				(SELECT id FROM new_match),
				($5)
		)
		RETURNING tournament_id;
	`

	var id int64
	if err := tx.QueryRow(ctx, sql, team1Id, team2Id, court_number, tournamentId, round_number).Scan(&id); err != nil {
		return fmt.Errorf("error while creating match: %w", err)
	}

	return nil
}

func queryCreateTournament(ctx context.Context, tx pgx.Tx, tournamentDate time.Time, tournamentType string) (int64, error) {

	sql := `
    INSERT INTO tournament (tournament_date, tournament_type_id)
    VALUES ($1, (SELECT id FROM tournament_type WHERE name = $2))
    RETURNING id;`

	var id int64
	if err := tx.QueryRow(ctx, sql, tournamentDate, tournamentType).Scan(&id); err != nil {
		return -1, fmt.Errorf("error while creating tournament: %w", err)
	}

	return id, nil
}

func CreateTournament(ctx context.Context, conn *pgxpool.Pool, t *tournament.Tournament) error {
	log.Print("creating tournament...")
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("msg rolling back transaction: %v", err)
		}
	}()

	tournamentType, err := tournament.TournamentTypeToString((*t).GetTournamentType())
	if err != nil {
		return fmt.Errorf("error converting tournament type to string: %w", err)
	}

	tournamentId, err := queryCreateTournament(ctx, tx, (*t).GetDateStart(), tournamentType)
	if err != nil {
		return err
	}

	teamIds := make(map[tournament.Team]int64)
	for _, team := range (*t).GetTeams() {
		teamId, err := queryInsertTeam(ctx, tx, team)
		if err != nil {
			return err
		}
		teamIds[team] = teamId
	}

	rounds := (*t).GetRounds()
	for roundIndex, round := range rounds {
		for _, match := range round {
			team1Id := teamIds[match.TeamA]
			team2Id := teamIds[match.TeamB]
			err := queryCreateMatch(ctx, tx, roundIndex, tournamentId, team1Id, team2Id, match.CourtId)
			if err != nil {
				return err
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}
