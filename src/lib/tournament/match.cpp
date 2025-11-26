#include "match.hpp"

const std::optional<const Team::Team*> Tournament::Match::get_team_1() const {
    return this->team_1;
}

const std::optional<const Team::Team*> Tournament::Match::get_team_2() const {
    return this->team_2;
}

std::ostream& Tournament::operator<<(std::ostream& os, const Match& match) {

    const auto& team_1 = match.team_1;
    const auto& team_2 = match.team_2;

    os << "Match(";
    if (team_1 && team_2) {
        os << *team_1.value() << ", " << *team_2.value() << ")";
    } else {
        os << "TBD";
    }

    os << ")";
    return os;
}