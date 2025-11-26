from dataclasses import dataclass
from functools import reduce
from math import ceil
from typing import Dict, FrozenSet, List, Optional, Set, Tuple


def get_matches_per_team(
    teams_number: int, total_rounds: int, available_courts: int
) -> Tuple[int, float, int]:
    matches_per_team = total_rounds

    while matches_per_team > 0:
        total_participations: int = teams_number * matches_per_team
        total_matches: float = total_participations / 2

        if total_matches % 1 == 0:
            matches_per_turn: float = total_matches / total_rounds

            if matches_per_turn <= available_courts:
                return int(total_matches), matches_per_turn, matches_per_team

        matches_per_team -= 1

    return 0, 0, 0


print("Total matches, Matches per turn, Matches per team")
print(get_matches_per_team(5, 5, 3))


def make_matches(teams: List[int], matches_per_team: int) -> Set[FrozenSet[int]]:
    n = len(teams)
    k = matches_per_team

    if not ((n * k) % 2 == 0 and n > k):
        return set()

    if k % 2 == 0:
        return k_regular_even(teams, k)

    return k_regular_odd(teams, k)


def k_regular_even(teams: List[int], k: int) -> Set[FrozenSet[int]]:
    n = len(teams)

    res: Set[FrozenSet[int]] = set()

    for i in range(0, n):
        for j in range(i - 1, i - 1 - k // 2, -1):
            res.add(frozenset({teams[j % n], teams[i]}))

        for j in range(i + 1, i + 1 + k // 2):
            res.add(frozenset({teams[j % n], teams[i]}))

    return res


def k_regular_odd(teams: List[int], k: int) -> Set[FrozenSet[int]]:
    n = len(teams)
    res = k_regular_even(teams, k - 1)

    for i in range(0, n):
        res.add(frozenset({teams[i], teams[(i + n // 2) % n]}))

    return res


print(make_matches([_ for _ in range(0, 5)], 4))


@dataclass
class Graph:
    nodes: Dict[int, Set[int]]

    def empty(self) -> bool:
        return all(len(v) == 0 for _, v in self.nodes.items())

    @staticmethod
    def from_set(edges: Set[FrozenSet[int]]) -> "Graph":
        n = len(reduce(lambda x, y: x | y, edges, set()))

        res = {i: set() for i in range(0, n)}

        for e in edges:
            assert len(e) == 2
            n1, n2 = e

            res[n1].add(n2)
            res[n2].add(n1)

        return Graph(res)


def make_turns(
    matches: Graph,
    matches_per_turn: float,
    turns: int,
) -> List[List[Tuple[int, int]]]:
    res: List[List[Tuple[int, int]]] = []
    n = ceil(matches_per_turn)
    playing_teams: Set[int] = set()

    turn_matches = []

    print("Graph" + f"{matches}")

    while len(res) < turns:
        for i in matches.nodes.keys():
            possible_matches = matches.nodes[i] - playing_teams
            if i not in playing_teams and len(possible_matches) > 0:
                player_1 = i
                player_2 = possible_matches.pop()
                playing_teams.add(player_1)
                playing_teams.add(player_2)
                turn_matches.append((player_1, player_2))
                matches.nodes[player_1].remove(player_2)
                matches.nodes[player_2].remove(player_1)
            if len(turn_matches) == n or (len(res) == turns - 1 and matches.empty()):
                res.append(turn_matches)
                playing_teams = set()
                turn_matches = []
        print("\n")
        print(res)
        print(matches)
        print(turn_matches)

    return res


print("\nTotal matches per turn\n")
print(make_turns(Graph.from_set(make_matches([_ for _ in range(0, 5)], 4)), 2, 5))

g = Graph.from_set(make_matches([_ for _ in range(0, 17)], 2))
tournament = make_turns(g, 4.25, 4)
print(tournament)


def validate_turns(
    graph: Graph, turns_matches: List[List[Tuple[int, int]]]
) -> Optional[str]:
    err = None
    overall_tournament: Dict[int, int] = {k: 0 for k in graph.nodes}
    for t in turns_matches:
        teams = set()
        for a, b in t:
            if a in teams:
                return f"Player {a} plays twice during turn {t}"
            if b in teams:
                return f"Player {b} plays twice during turn {t}"

            teams.add(a)
            teams.add(b)
            overall_tournament[a] += 1
            overall_tournament[b] += 1

    err = (
        "There is at least a team that plays less than the others"
        if len(set(v for _, v in overall_tournament.items())) > 1
        else None
    )

    return err


tournament_not_valid = validate_turns(g, tournament)
print(tournament_not_valid if tournament_not_valid else "Tournament valid")
