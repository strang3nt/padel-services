export interface Person {
  id: string;
}

export interface Team {
  person1: Person;
  person2: Person;
  gender: number;
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
}

export interface Tournaments {
  date: string;
  tournaments: TournamentData[];
}
