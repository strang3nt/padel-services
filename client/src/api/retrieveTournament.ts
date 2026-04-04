export default function retrieveTournament(
  bearerToken: string,
  date: Date,
): Promise<Response> {
  return fetch(`/api/tournaments?date=${date.toISOString()}`, {
    headers: { Authorization: `Bearer ${bearerToken}` },
  });
}
