package database

import (
	"context"
	"fmt"
	"time"

	"github.com/strang3nt/padel-services/pkg/tournament"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type tournamentNameType struct {
	TournamentId   int64
	TournamentType string
}

const tournamentsByDate = `
SELECT tournament.id, tournament_type.name
FROM tournament
JOIN tournament_type ON tournament.tournament_type_id=tournament_type.id
WHERE tournament.tournament_date::date = $1::date
`

type match struct {
	RoundNumber int
	Team1Id     int64
	Team2Id     int64
	CourtNumber int
}

const matchesByTournamentId = `
SELECT round_number, team1_id, team2_id, court_number
FROM match
JOIN round_tournament ON match.id=round_tournament.match_id
WHERE round_tournament.tournament_id=$1
ORDER BY round_number
`

type team struct {
	TeamId  int64
	Person1 string
	Person2 string
	Gender  string
}

const teamsByTournamentId = `
SELECT team.id, p1.name, p2.name, gender.name
FROM team
JOIN person p1 ON team.person1_id=p1.id
JOIN person p2 ON team.person2_id=p2.id
JOIN gender ON team.gender_id=gender.id
WHERE team.id IN (
	SELECT team1_id
	FROM match
	JOIN round_tournament ON match.id=round_tournament.match_id
	WHERE round_tournament.tournament_id=$1
	UNION
	SELECT team2_id
	FROM match
	JOIN round_tournament ON match.id=round_tournament.match_id
	WHERE round_tournament.tournament_id=$1
)
`

func GetTournamentsByDate(ctx context.Context, conn *pgxpool.Pool, tournamentDate time.Time) ([]tournament.TournamentData, error) {

	var tournaments []tournament.TournamentData
	rows, err := conn.Query(ctx, tournamentsByDate, tournamentDate)
	if err != nil {
		return tournaments, fmt.Errorf("query error: %w", err)
	}
	tournamentIds, err := pgx.CollectRows(rows, pgx.RowToStructByPos[tournamentNameType])
	if err != nil {
		return tournaments, fmt.Errorf("scan error: %w", err)
	}

	for _, id := range tournamentIds {
		rows, err := conn.Query(ctx, matchesByTournamentId, id.TournamentId)
		if err != nil {
			return tournaments, fmt.Errorf("query error: %w", err)
		}
		matches, err := pgx.CollectRows(rows, pgx.RowToStructByPos[match])
		if err != nil {
			return tournaments, fmt.Errorf("collectRows error: %v", err)
		}

		rows, err = conn.Query(ctx, teamsByTournamentId, id.TournamentId)
		if err != nil {
			return tournaments, fmt.Errorf("query error: %w", err)
		}
		teams, err := pgx.CollectRows(rows, pgx.RowToStructByPos[team])
		if err != nil {
			return tournaments, fmt.Errorf("collectRows error: %v", err)
		}

		tournaments = append(tournaments, buildTournamentData(id.TournamentType, tournamentDate, matches, teams))

	}

	return tournaments, nil

}

func buildTournamentData(
	tournamentName string,
	startDate time.Time,
	matches []match,
	teams []team) tournament.TournamentData {

	teamsMap := make(map[int64]*tournament.Team)
	matchesMap := make(map[int][]struct {
		team1        int64
		team2        int64
		court_number int
	})

	for _, m := range matches {
		curr, ok := matchesMap[m.RoundNumber]
		if !ok {
			matchesMap[m.RoundNumber] = []struct {
				team1        int64
				team2        int64
				court_number int
			}{
				{m.Team1Id, m.Team2Id, m.CourtNumber},
			}
		} else {
			curr = append(curr, struct {
				team1        int64
				team2        int64
				court_number int
			}{m.Team1Id, m.Team2Id, m.CourtNumber})
			matchesMap[m.RoundNumber] = curr
		}
	}

	teamsResult := make([]tournament.Team, 0)

	for _, t := range teams {
		person1 := tournament.Person{Id: t.Person1}
		person2 := tournament.Person{Id: t.Person2}

		teamsResult = append(teamsResult, tournament.Team{Person_1: person1, Person_2: person2, TeamGender: tournament.GenderFromString(t.Gender)})

		teamsMap[t.TeamId] = &teamsResult[len(teamsResult)-1]

	}

	rounds := make([]tournament.Round, len(matchesMap))

	for k, v := range matchesMap {
		round := make([]tournament.Match, 0)

		for _, t := range v {
			team1 := teamsMap[t.team1]
			team2 := teamsMap[t.team2]
			round = append(round, tournament.Match{
				TeamA:       *team1,
				TeamB:       *team2,
				MatchStatus: tournament.MatchScheduled,
				CourtId:     t.court_number})
		}
		rounds[k] = tournament.Round{Matches: round}
	}

	return tournament.MakeTournamentData(tournamentName, startDate, teamsResult, rounds)
}
