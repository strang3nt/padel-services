#pragma once
#include "person.hpp"
#include <iostream>

namespace Team {

enum class TeamGender {
    MALE,
    FEMALE,
    MIXED

};

enum class TeamLevel {
    BEGINNER,
    BEGINNER_TO_INTERMEDIATE,
    INTERMEDIATE,
    INTERMEDIATE_TO_ADVANCED,
    ADVANCED
};

class Team {
  private:
    Person first;
    Person second;

  public:
    const Person get_person_1() const;

    const Person get_person_2() const;

    Team(Person first, Person second) : first(first), second(second) {}

    double get_avg_age();
    TeamLevel get_team_level();
    TeamGender get_team_gender();

    bool operator==(const Team& other) const {
        return this->first == other.first && this->second == other.second;
    }
    bool operator!=(const Team& other) const { return !(*this == other); }

    friend std::ostream& operator<<(std::ostream& os, const Team& team);
};

std::ostream& operator<<(std::ostream& os, const Team& team);

} // namespace Team
