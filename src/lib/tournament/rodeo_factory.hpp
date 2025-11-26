#pragma once

#include "rodeo.hpp"
#include "tournament_factory.hpp"
#include <map>
#include <optional>
#include <set>

namespace Tournament {

using Edge = std::pair<int, int>;
using Matching = std::set<Edge>;
using AdjacencyMap = std::map<int, std::set<int>>;
using Matchings = std::vector<Matching>;

struct Graph {
    AdjacencyMap nodes;

    Graph(AdjacencyMap nodes_map) : nodes(std::move(nodes_map)) {}
    Graph() = default;

    /**
     * @brief Checks if the graph is "empty" (all nodes have no neighbors).
     * @return true if all adjacency sets are empty, false otherwise.
     */
    bool empty() const;

    /**
     * @brief Static method to construct a Graph from a set of unique edges.
     * * @param edges The set of unique matches (canonical pairs: {n1, n2} where n1 < n2).
     * @return A Graph object.
     */
    static Graph from_set(const Matching& edges);
};

/**
 * @todo write docs
 */
class RodeoFactory : public TournamentFactory {

  public:
    int total_rounds;
    int available_courts;

    RodeoFactory(int total_rounds, int available_courts)
        : TournamentFactory(), total_rounds(total_rounds), available_courts(available_courts) {}
    const Rodeo* make_tournament(const std::vector<const Team::Team*> teams) const;

    static std::optional<std::string> validate_rodeo(const Rodeo& matchings);

  private:
    /**
     * @brief Extracts a number of matchings
     *
     * @param graph The graph of remaining potential matches (will be modified).
     * @param avg_matching_size The ideal number of matches per turn (used to find max courts).
     * @param total_matchings The total number of turns to generate.
     * @return A list of turns, where each turn is a list of matches (pairs).
     */
    Matchings make_matchings(Graph& graph, double avg_matching_size, uint total_matchings) const;

    /**
     * @brief Main function to generate k-regular edges based on parity of k.
     *
     * @param nodes The list of nodes.
     * @param k The degree of each node.
     * @return A set of edges.
     */
    Matching make_edges(const std::vector<int>& nodes, int k) const;

    Matching k_regular_odd(const std::vector<int>& nodes, int k) const;

    Matching k_regular_even(const std::vector<int>& nodes, int k) const;

    /**
     * @brief Calculates the maximum number of matches per team, total matches,
     * and matches per round given constraints.
     *
     * @param teams_number The total number of teams participating.
     * @param total_rounds The maximum number of rounds available.
     * @param available_courts The number of courts available per round.
     * @return A tuple containing (total_matches, matches_per_turn, matches_per_team).
     */
    std::tuple<int, double, int> get_matches_per_team(int teams_number, int total_rounds,
                                                      int available_courts) const;

    void add_canonical_edge(Matching& res, int node_a, int node_b) const;

    Matchings make_matchings_brute_force(Matching& initial_edges, double avg_matching_size,
                                         uint total_matchings) const;
};

} // namespace Tournament
