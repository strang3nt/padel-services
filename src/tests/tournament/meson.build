

tournamenttest_sources = [
    'rodeo_factory_test.cpp'
]

tournamenttest_deps = [

]

tournamenttest_lib_deps = [
    tournament_lib
]

tournamenttest = executable('tournamenttest', 
    [tournamenttest_sources, test_main], 
    include_directories : libs_includes,
    dependencies : [tournamenttest_deps,test_deps], 
    link_with : [tournamenttest_lib_deps, test_dep_libs])

test('Tournament test', tournamenttest)

all_test_deps += tournamenttest_deps
all_test_dep_libs += tournamenttest_lib_deps
