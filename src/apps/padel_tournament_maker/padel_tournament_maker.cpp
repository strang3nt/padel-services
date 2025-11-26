#include <cstdlib>
#include <rodeo_factory.hpp>
#include <rodeo_maker.hpp>
#include <stdio.h>
#include <string>
#include <tgbot/tgbot.h>
#include <users_store.hpp>

const std::string TELEGRAM_BOT_API = std::getenv("TELEGRAM_BOT_API");

int main() {

    auto users_store = CsvUsersStore();

    TgBot::Bot bot(TELEGRAM_BOT_API);
    bot.getEvents().onCommand("start", [&bot, &users_store](TgBot::Message::Ptr message) {

        if (users_store.is_token_available(message->text)) {
            users_store.pair_token(message->chat->id);
            bot.getApi().sendMessage(message->chat->id, "Chat authorized");
        }
        else {
            bot.getApi().sendMessage(message->chat->id, "Chat not authorized");

        }


    });
    bot.getEvents().onAnyMessage([&bot, &users_store](TgBot::Message::Ptr message) {
        if (StringTools::startsWith(message->text, "/start")) {
            return;
        }

        if (!users_store.is_chat_authorized(message->chat->id)) {
            bot.getApi().sendMessage(message->chat->id, "Chat not authorized");
            return;
        }

        TgBot::Document::Ptr rounds = message->document;

        const TgBot::File::Ptr file_ptr = bot.getApi().getFile(rounds->fileId);
        const std::string file_str = bot.getApi().downloadFile(file_ptr->filePath);

        const Tournament::Rodeo* rodeo = RodeoMaker::make_from_csv(file_str);
        TgBot::InputFile::Ptr output_file = RodeoMaker::rodeo_to_csv(*rodeo);

        bot.getApi().sendDocument(message->chat->id, output_file);
    });
    try {
        printf("Bot username: %s\n", bot.getApi().getMe()->username.c_str());
        TgBot::TgLongPoll longPoll(bot);
        while (true) {
            printf("Long poll started\n");
            longPoll.start();
        }
    } catch (TgBot::TgException& e) {
        printf("error: %s\n", e.what());
    }
    return 0;
}
