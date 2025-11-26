#include <cstdio>
#include <person.hpp>
#include <rodeo_factory.hpp>
#include <rodeo_maker.hpp>
#include <sstream>
#include <string>
#include <team.hpp>

const Tournament::Rodeo* RodeoMaker::make_from_csv(const std::string& data) {

    std::istringstream f(data);
    std::string line;

    std::vector<const Team::Team*> teams;

    std::getline(f, line);
    int available_courts = std::stoi(line);
    std::getline(f, line);
    int turns = std::stoi(line);

    while (std::getline(f, line)) {
        std::vector<std::string> row;
        std::stringstream ss(line);
        std::string cell;

        while (std::getline(ss, cell, ',')) {
            row.push_back(cell);
        }

        Team::Person person_1{row[0]};
        Team::Person person_2{row[1]};

        teams.push_back(new Team::Team(person_1, person_2));
    }

    const Tournament::Rodeo* rodeo =
        Tournament::RodeoFactory(turns, available_courts).make_tournament(teams);

    return rodeo;
}
TgBot::InputFile::Ptr RodeoMaker::rodeo_to_csv(const Tournament::Rodeo& rodeo) {
    const auto teams = rodeo.get_teams();
    const auto turns = rodeo.get_turns();
    std::ofstream ss("test.txt");

    int i = 1;
    for (const auto& t : turns) {
        ss << "Round " << i << ",";
        int match = 1;
        for (const auto& m : t) {

            const auto team_1 = m.get_team_1().value();
            const auto team_2 = m.get_team_2().value();

            ss << "Match " << match << ",";
            ss << team_1->get_person_1().name << " - " << team_1->get_person_2().name << ",";
            ss << team_2->get_person_1().name << " - " << team_2->get_person_2().name << ",";
            match += 1;
        }
        ss << std::endl;
        i += 1;
    }

    ss.close();

    return TgBot::InputFile::fromFile("test.txt", "text/csv");
}
