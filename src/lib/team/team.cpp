#include "team.hpp"

double Team::Team::get_avg_age() { return -1.0; }

Team::TeamLevel Team::Team::get_team_level() { return TeamLevel::BEGINNER; }

Team::TeamGender Team::Team::get_team_gender() { return TeamGender::MALE; }

const Team::Person Team::Team::get_person_1() const { return this->first; }

const Team::Person Team::Team::get_person_2() const { return this->second; }

std::ostream& Team::operator<<(std::ostream& os, const Team& team) {
    os << "Team(" << team.first << ", " << team.second << ")";
    return os;
}