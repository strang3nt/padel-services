tournament_includes = include_directories('.')

tournament_sources = [
  'rodeo.cpp',
  'tournament.cpp',
  'match.cpp',
  'rodeo_factory.cpp',
  'tournament_factory.cpp',
]

tournament_lib = library('tournament',
  tournament_sources,
  include_directories: [team_includes],
  link_with: team_lib
)

