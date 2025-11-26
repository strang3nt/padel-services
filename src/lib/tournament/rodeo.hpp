#pragma once

#include "tournament.hpp"

namespace Tournament {

class Rodeo : public Tournament {

  private:
  public:
    Rodeo(const std::vector<const Team::Team*> teams, const std::vector<Turn> turns)
        : Tournament(teams, turns) {}

    ~Rodeo() {}
};

} // namespace Tournament
