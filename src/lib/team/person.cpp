#include "person.hpp"

std::ostream& Team::operator<<(std::ostream& os, const Person& person) {
    std::string gender = "male";

    switch (person.gender) {
    case Gender::FEMALE:
        gender = "female";
        break;
    default:
        break;
    }

    std::string level = "beginner";

    switch (person.level) {
    case Level::ADVANCED:
        level = "advanced";
        break;
    case Level::MEDIUM:
        level = "medium";
        break;
    default:
        break;
    }
    os << "Person(" << person.name << ", " << person.age << ", " << gender << ", " << level << ")";
    return os;
}
