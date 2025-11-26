#pragma once

#include <string>

class UsersStore {

public:
    virtual bool is_chat_authorized(int64_t chat_id) const = 0;
    virtual bool pair_token(int64_t chat_id) = 0;
    virtual bool is_token_available(std::string& token) const = 0;

};

class CsvUsersStore: public UsersStore {

public:
    bool is_chat_authorized(int64_t chat_id) const override;
    bool pair_token(int64_t chat_id) override;
    bool is_token_available(std::string& token) const override;
};