#pragma once

#include "../team/team.hpp"
#include <iostream>
#include <optional>

namespace Tournament {

enum class MatchStatus { OVER, ONGOING, NOT_STARTED, TBD };

struct MatchResult {
    int games_team_1{0};
    int games_team_2{0};

    bool operator==(const MatchResult& other) const {
        return this->games_team_1 == other.games_team_1 && this->games_team_2 == other.games_team_2;
    }

    bool operator!=(const MatchResult& other) const { return !(*this == other); }
};

class Match {
  private:
    std::optional<const Team::Team*> team_1;
    std::optional<const Team::Team*> team_2;
    MatchResult result;
    MatchStatus status;

  public:
    Match()
        : team_1(std::nullopt), team_2(std::nullopt), result(MatchResult{}),
          status(MatchStatus::TBD) {}

    Match(const Team::Team* team_1, const Team::Team* team_2)
        : team_1(team_1), team_2(team_2), status(MatchStatus::NOT_STARTED) {}

    void set_teams(const Team::Team* team_1, const Team::Team* team_2);
    void set_result(MatchResult result);

    bool operator==(const Match& other) const {
        return *this->team_1 == *other.team_1 && *this->team_2 == other.team_2 &&
               this->result == other.result && this->status == other.status;
    }

    bool operator!=(const Match& other) const { return *this == other; }

    friend std::ostream& operator<<(std::ostream& os, const Match& match);

    const std::optional<const Team::Team*> get_team_1() const;
    const std::optional<const Team::Team*> get_team_2() const;

    MatchStatus get_status() const;

    MatchResult get_result() const;
};

std::ostream& operator<<(std::ostream& os, const Match& match);

} // namespace Tournament
