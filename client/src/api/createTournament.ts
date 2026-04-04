import { Team, TournamentType } from "@/api/tournament";

export default function createTournament(
  bearerToken: string,
  eventName: string,
  tournamentType: TournamentType,
  dateStart: Date,
  roundsNumber: number,
  availableCourts: number,
  teams: Team[],
): Promise<Response> {
  return fetch(
    `/api/create-tournament?eventName=${eventName}&tournamentType=${tournamentType}&dateStart=${dateStart.toISOString()}&totalRounds=${roundsNumber}&availableCourts=${availableCourts}`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${bearerToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(teams),
    },
  );
}
