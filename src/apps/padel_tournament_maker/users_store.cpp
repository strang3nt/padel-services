#include <users_store.hpp>

bool CsvUsersStore::is_chat_authorized(int64_t chat_id) const {
    return true;
}

bool CsvUsersStore::pair_token(int64_t chat_id) {
    return true;
}

bool CsvUsersStore::is_token_available(std::string& token) const {
    return true;
};