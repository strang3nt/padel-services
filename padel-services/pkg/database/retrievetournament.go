package database

import (
	"context"
	"fmt"
	"padelservices/pkg/tournament"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

const tournamentsByDateAndType = `
SELECT tournament.tournament_id
FROM tournament
WHERE CURRENT_DATE(tournament.tournament_date)=CURRENT_DATE($1) AND tournament_type.name=$2
JOIN tournament_type ON tournament.tournament_type_id=tournament_type.id
`

type match struct {
	RoundNumber int
	Team1Id     int64
	Team2Id     int64
}

const matchesByTournamentId = `
SELECT round_number, team1_id, team2_id
FROM match, round_tournament
WHERE round_tournament.tournament_id=$1
JOIN round_tournament ON match.id=round_tournament.match_id
ORDER BY round_number
`

type team struct {
	teamId  int64
	person1 string
	person2 string
	gender  string
}

const teamsByTournamentId = `
SELECT team.id, p1.name, p2.name, gender.name
FROM team
WHERE team.id IN (
	SELECT team1_id
	FROM match
	WHERE round_tournament.tournament_id=$1
	JOIN round_tournament ON match.id=round_tournament.match_id
	UNION
	SELECT team2_id
	FROM match
	WHERE round_tournament.tournament_id=$1
	JOIN round_tournament ON match.id=round_tournament.match_id
)
JOIN person p1 ON team.person1_id=p1.id
JOIN person p2 ON team.person2_id=p2.id
JOIN gender ON team.gender_id=gender.id
`

func GetTournamentsByDate(ctx context.Context, conn *pgxpool.Pool, tournamentDate time.Time, tournamentType tournament.TournamentType) ([]tournament.TournamentData, error) {

	var tournaments []tournament.TournamentData
	tournamentTypeStr, err := tournament.TournamentTypeToString(tournamentType)
	if err != nil {
		return tournaments, fmt.Errorf("error while processing tournament type: %w", err)
	}

	rows, _ := conn.Query(ctx, tournamentsByDateAndType, tournamentDate, tournamentTypeStr)
	tournamentIds, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		return tournaments, fmt.Errorf("Scan error: %w", err)
	}

	for _, id := range tournamentIds {
		rows, err := conn.Query(ctx, matchesByTournamentId, id)
		matches, err := pgx.CollectRows(rows, pgx.RowToStructByPos[match])
		if err != nil {
			return tournaments, fmt.Errorf("CollectRows error: %v", err)
		}

		rows, err = conn.Query(ctx, teamsByTournamentId, id)
		teams, err := pgx.CollectRows(rows, pgx.RowToStructByPos[team])
		if err != nil {
			return tournaments, fmt.Errorf("CollectRows error: %v", err)
		}

		tournaments = append(tournaments, buildTournamentData(tournamentTypeStr, tournamentDate, matches, teams))

	}

	return tournaments, nil

}

func buildTournamentData(
	tournamentName string,
	startDate time.Time,
	matches []match,
	teams []team) tournament.TournamentData {
	var teamsMap map[int64]*tournament.Team

	var matchesMap map[int][]struct {
		team1 int64
		team2 int64
	}

	for _, m := range matches {
		curr, ok := matchesMap[m.RoundNumber]
		if !ok {
			matchesMap[m.RoundNumber] = []struct {
				team1 int64
				team2 int64
			}{
				{m.Team1Id, m.Team2Id},
			}
		} else {
			curr = append(curr, struct {
				team1 int64
				team2 int64
			}{m.Team1Id, m.Team2Id})
			matchesMap[m.RoundNumber] = curr
		}
	}

	teamsResult := make([]tournament.Team, 0)

	for _, t := range teams {
		person1 := tournament.Person{Id: t.person1}
		person2 := tournament.Person{Id: t.person2}

		teamsResult = append(teamsResult, tournament.Team{Person_1: person1, Person_2: person2, TeamGender: tournament.GenderFromString(t.gender)})

		teamsMap[t.teamId] = &teamsResult[len(teamsResult)-1]

	}

	rounds := make([]tournament.Round, len(matchesMap))

	for k, v := range matchesMap {
		round := make([]tournament.Match, 0)

		for _, t := range v {
			team1 := teamsMap[t.team1]
			team2 := teamsMap[t.team2]
			round = append(round, tournament.Match{TeamA: *team1, TeamB: *team2, MatchStatus: tournament.MatchScheduled})
		}
		rounds[k] = round
	}

	return tournament.MakeTournamentData(tournamentName, startDate, teamsResult, rounds)
}
