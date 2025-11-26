#pragma once
#include <iostream>
#include <string>

namespace Team {

enum class Gender { MALE, FEMALE };

enum class Level { BEGINNER, MEDIUM, ADVANCED };

struct Person {
    std::string name;
    int age{35};
    Gender gender{Gender::MALE};
    Level level{Level::BEGINNER};

    bool operator==(const Person& other) const {
        return (this->age == other.age) && (this->gender == other.gender) &&
               (this->level == other.level) && (this->name == other.name);
    }
    bool operator!=(const Person& other) const { return !(*this == other); }

    friend std::ostream& operator<<(std::ostream& os, const Person& person);
};

std::ostream& operator<<(std::ostream& os, const Person& person);

} // namespace Team