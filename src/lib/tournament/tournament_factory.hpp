#pragma once

#include "tournament.hpp"
#include <fstream>
#include <vector>

namespace Tournament {

class TournamentFactory {

  public:
    virtual const Tournament* make_tournament(const std::vector<const Team::Team*>) const = 0;
};

} // namespace Tournament
