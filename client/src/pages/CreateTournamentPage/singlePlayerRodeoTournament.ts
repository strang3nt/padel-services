export function getMatchesPerPerson(
  peopleNumber: number,
  totalRounds: number,
  availableCourts: number,
): [number, number, number] {
  if (peopleNumber === 0 || totalRounds === 0 || availableCourts === 0) {
    return [0, 0, 0];
  }

  const maxMatchesPerRound = Math.min(
    Math.floor(peopleNumber / 4),
    availableCourts,
  );
  const maxTotalMatches = maxMatchesPerRound * totalRounds;
  let k = Math.min(
    Math.floor((maxTotalMatches * 4) / peopleNumber),
    totalRounds,
  );

  while ((peopleNumber * k) % 4 != 0 && k > 0) {
    k -= 1;
  }

  const totalMatches = (k * peopleNumber) / 4;
  const matchesPerRound = Math.ceil(totalMatches / totalRounds);

  return [totalMatches, matchesPerRound, k];
}
