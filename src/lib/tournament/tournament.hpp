#pragma once

#include "match.hpp"
#include <ranges>
#include <sstream>
#include <string>
#include <vector>

namespace Tournament {

using Turn = std::vector<Match>;

class Tournament {
  private:
    const std::vector<const Team::Team*> teams;
    const std::vector<Turn> turns;

  public:
    Tournament(const std::vector<const Team::Team*> teams, const std::vector<Turn> turns)
        : teams(teams), turns(turns) {}

    virtual ~Tournament() = 0;

    const std::vector<const Team::Team*> get_teams() const { return this->teams; }

    const std::vector<Turn> get_turns() const { return this->turns; }

    bool operator==(const Tournament& other) const {

        if (this->teams.size() != other.teams.size()) {
            return false;
        }

        if (this->turns.size() != other.turns.size()) {
            return false;
        }

        for (uint i = 0; i < this->teams.size(); ++i) {
            const Team::Team team_1 = *this->teams.at(i);
            const Team::Team team_2 = *other.teams.at(i);
            if (!(team_1 == team_2)) {
                return false;
            }
        }

        for (uint i = 0; i < this->turns.size(); ++i) {
            const auto turn_1 = this->turns.at(i);
            const auto turn_2 = other.turns.at(i);
            if (!(turn_1 == turn_2)) {
                return false;
            }
        }

        return true;
    }

    bool operator!=(const Tournament& other) const { return !(*this == other); }

    friend std::ostream& operator<<(std::ostream& os, const Tournament& tournament);
};

std::ostream& operator<<(std::ostream& os, const Tournament& tournament);

} // namespace Tournament
