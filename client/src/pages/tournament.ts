export interface Person {
  id: string;
}

export interface Team {
  person1: Person;
  person2: Person;
  gender: number;
}

export function genderToString(g: number): string {
  switch (g) {
    case 0:
      return "Male";
    case 1:
      return "Female";
    case 2:
      return "Mixed";
    default:
      throw Error("Team gender not recognized");
  }
}

export interface Match {
  teamA: Team;
  teamB: Team;
  matchStatus: number;
  courtId: number;
}

export interface Matches {
  matches: Match[];
}

export interface TournamentData {
  name: string;
  date: string;
  teams: Team[];
  rounds: Matches[];
  tournamentType: string;
}

export interface Tournaments {
  date: string;
  tournaments: TournamentData[];
}
