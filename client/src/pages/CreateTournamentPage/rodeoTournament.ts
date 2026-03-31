export function getMatchesPerTeam(
  teamsNumber: number,
  totalRounds: number,
  availableCourts: number,
): [number, number, number] {
  let matchesPerTeam = totalRounds;

  while (matchesPerTeam > 0) {
    const totalParticipations = teamsNumber * matchesPerTeam;
    const totalMatchesFloat = totalParticipations / 2.0;

    if (Number.isInteger(totalMatchesFloat)) {
      const matchesPerTurn = totalMatchesFloat / totalRounds;

      if (matchesPerTurn <= availableCourts && teamsNumber > matchesPerTeam) {
        return [totalMatchesFloat, matchesPerTurn, matchesPerTeam];
      }
    }

    matchesPerTeam -= 1;
  }

  return [0, 0.0, 0];
}
