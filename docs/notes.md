# Analysis of a Rodeo-tournament maker algorithm

## Introduction and motivation

> Padel is a game inspired from tennis, from which inherits
> most of its rules. The peculiarity is that the court is
> surrounded by see-through walls, where the ball can do its second
> bounce. The game is tipically played in doubles.

Rodeo is a type of tournament that is, among others, employed in Padel.
The general set of rules is that matches are quick (e.g. one-game matches),
teams all play the same number of matches and the team with most points wins
tournament.
Other variations have players switch teams, and the player with most points wins.
Some add a semi-final and final round in case of draw among participants.

This is the kind of Rodeo I am going to focus on:

- teams all play the same numbers of matches
- matches last one game
- winner is the team with the most points.

A draw is handled via a final match (if between two teams) or a round-robin
round.

I focus only on the first part of the tournament, more
specifically on the match-making problem.
In the remainder of this report I first describe the problem and then
show an algorithm that does the match-making.

## The Rodeo-match-making problem

More formally, these are the set of constraints I have to follow.

- Teams all play the same number of matches thoughout the tournament.
- Teams play at most a match per round, they may rest.
- Teams never play against each other more than once throughout the tournament.

> By "round" I mean a set of matches that are all played
> during the same timeframe. A round ends when all matches in the same
> round end.

The parameters I have to work with are:

- the teams $t$;
- the number of rounds $r$;
- the number of Padel courts available $p$.

The output is a list of lists of matches, each list assigned to a round.
The problem can be split in three subproblems.

1. Given the number of teams $t$, rounds $r$, padel courts $p$, I want to know
the maximum number of matches per team, $m$.
2. Given the matches per team I compute the pairings.
3. I assign the set of matches between the different rounds.

Let's dive in!

### Number of matches per team

This is just some math. Given the number of teams $t$, rounds $r$, courts available $p$,
compute the most matches per team.
The algorithm below tries to find via a brute force approach a solution.
First, it assumes there is a solution where each team plays once every round, that is, the
situation where there are at least $t / 2$ padel courts available.
Then, reduces iteratively by one the candidate matches per team, until it either finds a solution
or reaches 0. 0 means no solution has been found. This is a situation that can happen with
ill formed inputs, such as 0 courts available or an input that does not allow every team to
play at leas once.

In the while loop, the algorithm simply tries to find a total number of matches that
are even, and if that number leads to a number of matches per turn that is lower than $p$,
we found a solution.

> Note that we do not care in this step whether `matches_per_turn` is integer, we will deal
> with it later when actually filling the rounds.

The algorithm finally returns the total matches, matches per turn and matches per team.

```python
def get_matches_per_team(t, r, p):
    
    matches_per_team = r

    while matches_per_team > 0:
        total_participations = t * matches_per_team
        total_matches = total_participations / 2

        if total_matches % 1 == 0:
            matches_per_turn = total_matches / r

            if matches_per_turn <= p:
                return total_matches, matches_per_turn, matches_per_team

        matches_per_team -= 1

    return 0, 0, 0
```

### Matchmaking

Things are starting to get complicated. For each team, we need to find $m$ matches.
It easy to see a graph data structure where nodes are teams, and edges are the matches.
In fact, we are not building any graph, but a k-regular one!
We have a set of nodes (teams), and each node must have a degree $k = m$.

The preconditios for building a k-regular graph are that the cardinality is
even and that the number of nodes is larger than the number of edges, i.e.
$(n * k) \% 2 == 0\wedge n > k$. Both of which we satisfy, since $(n * k)$
is the number of total participations which must be even (each match counts as
two participations) and the number of teams is of course larger than the number
of matches per team, we would have otherwise a team playing multiple times against
another team.

we must distinguish between the case where k is even or odd.

```python
def make_matches(teams, matches_per_team):
    n = len(teams)
    k = matches_per_team

    if k % 2 == 0:
        return k_regular_even(teams, k)

    return k_regular_odd(teams, k)
```

The idea is the following, we arrange the nodes in a circle, and:

- if k is even, we select $k/2$ nodes to the left and right of each node;
- if k is odd, we do the same and for each node $n$ we draw an additional edge to
that is exactly on the opposite side of $n$ (in the imaginary circle).

```python
def k_regular_even(teams, k):
    n = len(teams)

    for i in range(0, n):
        for j in range(i - 1, i - 1 - k // 2, -1):
            res.add({teams[j % n], teams[i]})

        for j in range(i + 1, i + 1 + k // 2):
            res.add({teams[j % n], teams[i]})

    return res
```

```python
def k_regular_odd(teams: List[int], k: int) -> Set[FrozenSet[int]]:
    n = len(teams)
    res = k_regular_even(teams, k - 1)

    for i in range(0, n):
        res.add(frozenset({teams[i], teams[(i + n // 2) % n]}))

    return res
```

### Assigning matches to rounds

Again, this problem can be generalized to a graph problem. We have a set of
matches (edges), and teams (nodes) and we want to find a way to separate all the
matches into different rounds in such a way that each in each round no team plays
twice, that is, each round is a matching [^1] and we want to create as many matchings
as there are rounds. This is an edge-coloring problem [^2].
Given a k-regular graph, split the edges into $n$ sets of edges, where every
set contains non-adjacent edges (a set of non-adjacent edges is an indipendent set,
or a matching). I need to find $n$ matchings where $n$ is the number of rounds. and
the cardinality of the sets is the number of matches for each round.
Edge-colorign is a well-studied and a very hard problem, NP-hard to be precise.
Luckily, the input is a k-regular graph which makes the implementation fairly straight forward.
The intuition is: just pick one edge per vertex to build a round.
Until the last round, we have at least $ceil(m)$ matches that we can extract.

Sike! The algorithm above is wrong! It was the first implementation, and it had a problem:
it could lead to a situation while building the last round only adjacent edges remain!

I chose to use a brute-force solution. A better solution,
would not be polynomial anyways, unless I misunderstood the problem's assumptions.
Intuitively, a brute force algorithm would guess at each step in which round a
specific match would go. It builds a decision tree, and backtracks whenever it
realizes it took the wrong path. The stop condition is that all matches have been
assigned to a round, while satisfying the constraint of no team playing twice during
a round and rounds are less or equal than the previously calculated $m$ variable,
matches-per-round.
Doing some quick napkin math we find that the tree of choiches a brute-force algorithm
may go through is huge. The complexity is exponential and should be around
$K^E$ where $E$ is the total number of matches during the tournament (and the depth
of the recursion tree), and $K$ the number of rounds (where to place the match).

Below a backtracking algorithm. Let me remark the parallelism between the
graph problem and the actual tournament-making problem.

- **remaining_edges**, matches not yet scheduled.
- **buckets**, matches scheduled for each round.
- **used_nodes**, for each round, the teams already scheduled.
- **total_matchings**, the number of rounds.
- **matching_max_size**, the maximum number of matches per round.

```python
def find_matchings(
    remaining_edges,
    buckets,
    used_nodes,
    total_matchings,
    matching_max_size
):

    if len(remaining_edges) == 0:
        return buckets

    current_edge = remaining_edges.pop()
    p1, p2 = current_edge

    for i in range(0, total_matchings - 1):
        
        if (
            (p1 not in used_nodes[i]) and (p2 not in used_nodes[i]) 
            and (size(buckets[i]) < matching_max_size)
        ):

            buckets[i].add(current_edge)
            used_nodes[i].add(p1)
            used_nodes[i].add(p2)

            result = Solve(
                remaining_edges,
                buckets,
                used_nodes,
                total_matchings,
                matching_max_size
            )
            if result:
                return result

            buckets[i].remove(current_edge)
            used_nodes[i].remove(p1)
            used_nodes[i].remove(p2)

    remaining_edges.push(current_edge)
    return false
```

This is an exponential algorithm, does it work in practice? Let me assume that
I am a magician and that building a state in the decision tree
takes only 1 CPU cycle, and that I have a 4 GHz CPU, thus roughly 4*10^9 operations
per second. I want to build a tournament with 15 teams, 5 courts
available and 8 rounds, this is a real situation. Thus, $K = 8$ and $E = 30$.
Doing some quick math, it should take $10$ billion years to compute the whole solution space.
Hopefully the algorithm converges to a path to the correct solution fast!

## Conclusion

This was hard! Tournament making (not only in the case of Rodeo tournaments) is a difficult
problem that require complex and often slow algorithms. We found that this particular tournament
was just an instance of an edge-coloring problem. There is already a lot of
literature that try to solve this problem on graphs that satisfy certain properties,
and provide algorithms. Still, the problem is classified to be NP-hard and thus
only exponential algorithms are known to be able to solve the general case.
I provided a backtracking algorithm that builds a solution.
Improvements to the algorithm could be made by adding strategies to prune the search
space earlier, for example by choosing the edges according to some order.
Higher level optimizations might also be interesting to add:

- I could apply a heuristic first, and only if it does not work, use the brute-force
algorithm.
- I could add some form of randomness and run instances of the algorithm in
parallel, and wait for the first one to finish.

[^1] <https://en.wikipedia.org/wiki/Matching_(graph_theory)>
[^2] <https://en.wikipedia.org/wiki/Edge_coloring>
