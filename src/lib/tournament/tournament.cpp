#include "tournament.hpp"
#include "../team/team.hpp"

Tournament::Tournament::~Tournament() {}

std::ostream& Tournament::operator<<(std::ostream& os, const Tournament& tournament) {
    os << "Tournament(Teams(";

    for (const auto t : tournament.teams) {
        os << *t << ", ";
    }

    os << "), Turns(";

    for (const auto& t : tournament.turns) {
        os << "Turn(";
        for (const auto& m : t) {
            os << m << ", ";
        }
        os << "), ";
    }

    os << ")";
    return os;
}
