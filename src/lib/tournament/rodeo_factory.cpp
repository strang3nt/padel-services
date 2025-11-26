#include "rodeo_factory.hpp"
#include "tournament.hpp"
#include <algorithm>
#include <cassert>
#include <cmath>
#include <iostream>
#include <memory>
#include <numeric>
#include <optional>
#include <set>
#include <sstream>
#include <stack>
#include <tuple>
#include <unordered_set>
#include <utility>
#include <vector>

const Tournament::Rodeo*
Tournament::RodeoFactory::make_tournament(const std::vector<const Team::Team*> teams) const {

    int n = teams.size();
    std::vector<int> nodes(n);

    std::iota(nodes.begin(), nodes.end(), 0);

    int total_matches;
    double matches_per_turn;
    int matches_per_team;

    std::tie(total_matches, matches_per_turn, matches_per_team) =
        this->get_matches_per_team(n, this->total_rounds, this->available_courts);

    Matching all_matches = this->make_edges(nodes, matches_per_team);


    Graph graph = Graph::from_set(all_matches);
    Matchings rounds = this->make_matchings(graph, matches_per_turn, this->total_rounds);

    if (!graph.empty()){
        rounds = this->make_matchings_brute_force(all_matches, matches_per_turn, this->total_rounds);
    }

    std::vector<Turn> turns;

    for (const auto& matching : rounds) {
        std::vector<Match> matches;
        for (const auto& edge : matching) {
            int e_1;
            int e_2;
            std::tie(e_1, e_2) = edge;
            Match m(teams.at(e_1), teams.at(e_2));
            matches.push_back(m);
        }
        turns.push_back(matches);
    }

    return new Rodeo(teams, turns);
}

std::tuple<int, double, int>
Tournament::RodeoFactory::get_matches_per_team(int teams_number, int total_rounds,
                                               int available_courts) const {

    int matches_per_team = total_rounds;

    while (matches_per_team > 0) {

        int total_participations = teams_number * matches_per_team;
        double total_matches = (double)total_participations / 2.0;

        if (std::floor(total_matches) == total_matches) {
            double matches_per_turn = total_matches / total_rounds;

            if (matches_per_turn <= available_courts) {

                return std::make_tuple((int)total_matches, matches_per_turn, matches_per_team);
            }
        }

        matches_per_team -= 1;
    }

    return std::make_tuple(0, 0.0, 0);
}

void Tournament::RodeoFactory::add_canonical_edge(Matching& res, int node_a, int node_b) const {
    if (node_a < node_b) {
        res.insert({node_a, node_b});
    } else {
        res.insert({node_b, node_a});
    }
}

Tournament::Matching Tournament::RodeoFactory::k_regular_even(const std::vector<int>& nodes,
                                                              int k) const {
    int n = nodes.size();
    Matching res;

    for (int i = 0; i < n; ++i) {
        for (int count = 1; count <= k / 2; ++count) {
            int j_index = i - count;
            int j_mod_n = (j_index % n + n) % n;

            add_canonical_edge(res, nodes[j_mod_n], nodes[i]);
        }

        for (int count = 1; count <= k / 2; ++count) {
            int j_index = i + count;
            int j_mod_n = j_index % n;

            add_canonical_edge(res, nodes[j_mod_n], nodes[i]);
        }
    }

    return res;
}

Tournament::Matching Tournament::RodeoFactory::k_regular_odd(const std::vector<int>& nodes,
                                                             int k) const {
    int n = nodes.size();

    Matching res = k_regular_even(nodes, k - 1);

    for (int i = 0; i < n; ++i) {
        int partner_index = (i + n / 2) % n;

        this->add_canonical_edge(res, nodes[i], nodes[partner_index]);
    }

    return res;
}

Tournament::Matching Tournament::RodeoFactory::make_edges(const std::vector<int>& nodes,
                                                          int k) const {
    int n = nodes.size();

    if (!(((n * k) % 2 == 0) && (n > k))) {
        return Matching();
    }

    if (k % 2 == 0) {
        return k_regular_even(nodes, k);
    }

    return k_regular_odd(nodes, k);
}

bool Tournament::Graph::empty() const {
    for (const auto& pair : nodes) {
        if (!pair.second.empty()) {
            return false;
        }
    }
    return true;
}

Tournament::Graph Tournament::Graph::from_set(const Matching& edges) {

    int max_node_id = -1;
    for (const auto& edge : edges) {
        max_node_id = std::max(max_node_id, std::max(edge.first, edge.second));
    }

    AdjacencyMap res;
    for (int i = 0; i <= max_node_id; ++i) {
        res[i] = std::set<int>();
    }

    for (const auto& edge : edges) {
        int n1 = edge.first;
        int n2 = edge.second;

        res[n1].insert(n2);
        res[n2].insert(n1);
    }

    return Graph(res);
}

Tournament::Matchings Tournament::RodeoFactory::make_matchings_brute_force(
    Matching& initial_edges, double avg_matching_size, uint total_matchings) const {

    avg_matching_size = std::ceil(avg_matching_size);
    struct State {
        Matchings buckets;
        Matching remaining_edges;
        std::vector<std::set<int>> buckets_used_nodes;
    };

    std::stack<std::unique_ptr<State>> stack;

    stack.push(std::make_unique<State>(std::vector<Matching>(total_matchings, Matching()),
                                       initial_edges, std::vector<std::set<int>>(total_matchings)));

    while (!stack.empty()) {

        std::unique_ptr<State> current_state_ptr = std::move(stack.top());
        stack.pop();
        State& current_state = *current_state_ptr;

        if (current_state.remaining_edges.empty()) {
            return current_state.buckets;
        }

        auto it = current_state.remaining_edges.begin();
        Edge current_edge = *it;
        int p1 = current_edge.first;
        int p2 = current_edge.second;

        for (uint i = 0; i < total_matchings; ++i) {
            const auto& players = current_state.buckets_used_nodes[i];
            const auto& edges_in_matching = current_state.buckets[i].size();
            if (players.find(p1) == players.end() && players.find(p2) == players.end() &&
                edges_in_matching < avg_matching_size) {

                std::unique_ptr<State> next_state_ptr = std::make_unique<State>(current_state);
                State& next_state = *next_state_ptr;
                next_state.buckets[i].insert(current_edge);

                next_state.buckets_used_nodes[i].insert(p1);
                next_state.buckets_used_nodes[i].insert(p2);

                int erased = next_state.remaining_edges.erase(current_edge);
                assert(erased == 1);
                stack.push(std::move(next_state_ptr));
            }
        }
    }

    return Matchings();
}

Tournament::Matchings Tournament::RodeoFactory::make_matchings(Graph& graph,
                                                               double avg_matching_size,
                                                               uint total_matchings) const {

    Matchings res;
    int max_matches_per_turn = (int)std::ceil(avg_matching_size);
    std::set<int> playing_teams;

    std::set<std::pair<int, int>> turn_matches;

    while (res.size() < total_matchings) {

        bool added_something = false;

        for (const auto& pair : graph.nodes) {
            int i = pair.first;

            std::set<int> possible_matches;
            std::set_difference(graph.nodes[i].begin(), graph.nodes[i].end(), playing_teams.begin(),
                                playing_teams.end(),
                                std::inserter(possible_matches, possible_matches.begin()));

            if (playing_teams.find(i) == playing_teams.end() && !possible_matches.empty()) {
                added_something = true;
                int player_2 = *possible_matches.begin();
                int player_1 = i;

                playing_teams.insert(player_1);
                playing_teams.insert(player_2);

                turn_matches.insert({player_1, player_2});

                graph.nodes[player_1].erase(player_2);
                graph.nodes[player_2].erase(player_1);
            }

            if (turn_matches.size() == (size_t)max_matches_per_turn ||
                (res.size() == total_matchings - 1 && graph.empty())) {

                res.push_back(turn_matches);
                playing_teams.clear();
                turn_matches.clear();
            }
        }

        if (!added_something) {
            break;
        }
    }

    return res;
}

std::optional<std::string> Tournament::RodeoFactory::validate_rodeo(const Rodeo& rodeo) {
    const auto& teams = rodeo.get_teams();
    std::map<const Team::Team*, int> team_to_int;

    int total_teams = teams.size();

    std::map<int, int> overall_tournament;
    for (int i = 0; i < total_teams; ++i) {
        team_to_int[teams[i]] = i;
        overall_tournament[i] = 0;
    }

    Matchings matchings;
    for (const auto& t : rodeo.get_turns()) {
        matchings.push_back(Matching());
        for (const auto& m : t) {
            const std::optional<const Team::Team*> team_1 = m.get_team_1();
            const std::optional<const Team::Team*> team_2 = m.get_team_2();

            if (!team_1 || !team_2) {
                std::stringstream ss;
                ss << "There is a match whose teams were not all set";
                return ss.str();
            }

            matchings.back().insert(
                {team_to_int.at(team_1.value()), team_to_int.at(team_2.value())});
        }
    }

    for (const auto& t : matchings) {

        std::set<int> teams_in_turn;

        for (const auto& match : t) {
            int a = match.first;
            int b = match.second;

            if (teams_in_turn.count(a)) {
                std::stringstream ss;
                ss << "Player " << a << " plays twice during turn (teams_in_turn)";
                return ss.str();
            }
            if (teams_in_turn.count(b)) {
                std::stringstream ss;
                ss << "Player " << b << " plays twice during turn (teams_in_turn)";
                return ss.str();
            }

            teams_in_turn.insert(a);
            teams_in_turn.insert(b);

            overall_tournament[a]++;
            overall_tournament[b]++;
        }
    }

    std::set<int> unique_match_counts;
    for (const auto& pair : overall_tournament) {
        unique_match_counts.insert(pair.second);
    }

    if (unique_match_counts.size() > 1) {
        return "There is at least a team that plays less than the others";
    }

    return std::nullopt;
}
