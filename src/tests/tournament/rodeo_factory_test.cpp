#include "gtest/gtest.h"
#include <string>
#include <team/person.hpp>
#include <team/team.hpp>
#include <tournament/match.hpp>
#include <tournament/rodeo.hpp>
#include <tournament/rodeo_factory.hpp>
#include <tournament/tournament.hpp>
#include <vector>

TEST(TournamentTests, rodeoFactory) {
    const Team::Team* team_0 = new Team::Team(Team::Person{"1"}, Team::Person{"2"});
    const Team::Team* team_1 = new Team::Team(Team::Person{"3"}, Team::Person{"4"});
    const Team::Team* team_2 = new Team::Team(Team::Person{"5"}, Team::Person{"6"});
    const Team::Team* team_3 = new Team::Team(Team::Person{"7"}, Team::Person{"8"});
    const Team::Team* team_4 = new Team::Team(Team::Person{"9"}, Team::Person{"10"});
    const std::vector<const Team::Team*> teams{team_0, team_1, team_2, team_3, team_4};

    const std::vector<Tournament::Match> turn_1{Tournament::Match(team_4, team_0),
                                                Tournament::Match(team_1, team_3)};
    const std::vector<Tournament::Match> turn_2{Tournament::Match(team_4, team_1),
                                                Tournament::Match(team_2, team_0)};
    const std::vector<Tournament::Match> turn_3{Tournament::Match(team_4, team_2),
                                                Tournament::Match(team_3, team_1)};
    const std::vector<Tournament::Match> turn_4{Tournament::Match(team_4, team_3),
                                                Tournament::Match(team_0, team_2)};
    const std::vector<Tournament::Match> turn_5{Tournament::Match(team_0, team_1),
                                                Tournament::Match(team_2, team_3)};

    const auto rodeo_tournament = Tournament::Rodeo(
        teams, std::vector<Tournament::Turn>{turn_1, turn_2, turn_3, turn_4, turn_5});

    const auto actual_rodeo_tournament = *Tournament::RodeoFactory(5, 2).make_tournament(teams);
    EXPECT_EQ(rodeo_tournament, actual_rodeo_tournament);
}

TEST(TournamentTests, rodeoValidationWrongTournament) {
    const Team::Team* team_0 = new Team::Team(Team::Person{"1"}, Team::Person{"2"});
    const Team::Team* team_1 = new Team::Team(Team::Person{"3"}, Team::Person{"4"});
    const Team::Team* team_2 = new Team::Team(Team::Person{"5"}, Team::Person{"6"});
    const Team::Team* team_3 = new Team::Team(Team::Person{"7"}, Team::Person{"8"});
    const Team::Team* team_4 = new Team::Team(Team::Person{"9"}, Team::Person{"10"});
    const std::vector<const Team::Team*> teams{team_0, team_1, team_2, team_3, team_4};

    const std::vector<Tournament::Match> turn_1{Tournament::Match(team_4, team_0),
                                                Tournament::Match(team_1, team_3)};
    const std::vector<Tournament::Match> turn_2{Tournament::Match(team_4, team_1),
                                                Tournament::Match(team_2, team_0)};
    const std::vector<Tournament::Match> turn_3{Tournament::Match(team_4, team_2),
                                                Tournament::Match(team_3, team_1)};
    const std::vector<Tournament::Match> turn_4{Tournament::Match(team_4, team_3),
                                                Tournament::Match(team_0, team_2)};
    const std::vector<Tournament::Match> turn_5{Tournament::Match(team_4, team_3),
                                                Tournament::Match(team_2, team_3)};

    const auto rodeo_tournament = Tournament::Rodeo(
        teams, std::vector<Tournament::Turn>{turn_1, turn_2, turn_3, turn_4, turn_5});

    ASSERT_TRUE(Tournament::RodeoFactory::validate_rodeo(rodeo_tournament));
}

TEST(TournamentTests, rodeoFactoryValidation) {
    const Team::Team* team_0 = new Team::Team(Team::Person{"1"}, Team::Person{"2"});
    const Team::Team* team_1 = new Team::Team(Team::Person{"3"}, Team::Person{"4"});
    const Team::Team* team_2 = new Team::Team(Team::Person{"5"}, Team::Person{"6"});
    const Team::Team* team_3 = new Team::Team(Team::Person{"7"}, Team::Person{"8"});
    const Team::Team* team_4 = new Team::Team(Team::Person{"9"}, Team::Person{"10"});
    const std::vector<const Team::Team*> teams{team_0, team_1, team_2, team_3, team_4};

    const auto rodeo_tournament = *Tournament::RodeoFactory(5, 2).make_tournament(teams);

    ASSERT_FALSE(Tournament::RodeoFactory::validate_rodeo(rodeo_tournament));
}