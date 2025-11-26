# Analysis of a Rodeo-tournament maker algorithm

In the (italian) Padel scene Rodeo tournaments are very popular.
Rodeo has no official set of rules, but in general matches are
quick (e.g. one-game matches), teams all play the same number of
matches and in general the team with most points/wins, wins the
tournament.
Interesting variations have players play with different, random
teams, and the player with most points wins. Others add a semi-final
and final round.

The kind of Rodeo I focus on is the following:

- teams all play the same numbers of matches
- matches last one game
- winner is the team with the most points.

A draw is handled via a final match (if between two teams) or round-robin
round.

Here we focus on the first part of the tournament, and more
specifically on the match-making problem.
In the remainder of this report I first describe the problem and then
show an algorithm that does the match-making for a Rodeo tournament.

## The Rodeo-match-making problem

These are the set of constraints I have to follow:

- teams all play the same number of matches thoughout the tournament;
- teams play at most a match per round (they might also rest);
- teams never play against the same team more than once.

Note that by "round" I mean a set of matches that are all played
during the same timeframe. A round ends when all matches in the same
round end.

There are also obvious constraints such as a team cannot play more
than one match during the same round.

The parameters I have to work with are:

- the teams;
- the number of rounds;
- the number of Padel courts available.

The output is a list of matches separated in different rounds.

I decided to split the problem in three subproblems:

1. given the number of teams, rounds, padel courts, I want to know
the maximum number of matches per team
2. given the matches per team (which I know from the previous step)
I compute the matches per team
3. once I have a set of matches, I split them between the different
rounds.

Let's dive in!

### Number of matches per team

Easy

### Matchmaking

K-regular graph! This is akin to the problem of building a k-regular graph,
we have a set of nodes (teams), and each node have a degree of exatly k. k is
the number of matches per team!

### Assigning matches to rounds

Given a k-regular graph, split the edges into n sets of edges, where every
set contains non-adjacent edges (a set of non-adjacent edges is an indipendent set,
or a matching). Thus I need to find n matchings where n is the number of rounds. and
the cardinality of the sets is less or equal to fields available and n is the number
of rounds.
This is akin to finding a n-edge-coloring of the graph, that is, the problem of
coloring a graph with n colors s.t. the colors are non adjacent. This is a well-studies
and very hard problem (NP-hard to be precise). Luckily for me, the input is a
k-regular graph which makes the implementation fairly straight forward. The intuition
is: just pick one edge per vertex to build a round. Until the last round, we have at least
ceil(m) matches that we can extract. 

