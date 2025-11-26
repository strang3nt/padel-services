#pragma once

#include <rodeo.hpp>
#include <rodeo_factory.hpp>
#include <string>
#include <tgbot/tgbot.h>

class RodeoMaker {
  public:
    static const Tournament::Rodeo* make_from_csv(const std::string& data);
    static TgBot::InputFile::Ptr rodeo_to_csv(const Tournament::Rodeo& rodeo);
};
