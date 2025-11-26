cmake = import('cmake')

tgbot_proj = cmake.subproject('tgbot-cpp')
tgbot_dep = tgbot_proj.dependency('TgBot')

executable(
  'telegram_bot',
  'padel_tournament_maker.cpp',
  'rodeo_maker.cpp',
  'users_store.cpp',
  link_with: [team_lib, tournament_lib],
  include_directories: [team_includes, tournament_includes],
  dependencies: [tgbot_dep]
)


