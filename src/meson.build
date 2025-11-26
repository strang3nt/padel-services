project('padel-tournament-maker', 'cpp',
  version : '0.1',
  default_options : {
    'warning_level': '3',
    'cpp_std': 'c++20',
    'c_std': 'c11',
    'default_library': 'static',
  }
)

gtest = subproject('gtest')

# add cppcheck inspector target
cppcheck = find_program('cppcheck', required : false)
if cppcheck.found()

  subprojects_dir = join_paths(meson.global_source_root(), 'subprojects', '**')

  cppcheck_args = [
    '--project=' + join_paths(meson.project_build_root(), 'compile_commands.json'),
    # '--clang',
    '--suppress=*:' + subprojects_dir,

    '--cppcheck-build-dir='+meson.current_build_dir(),
    '--std=c++20',
    '--enable=all',
    '--inconclusive'
  ]
  run_target('cppcheck',
    command : [cppcheck, cppcheck_args]
  )
endif


subdir('lib')
subdir('apps')
subdir('tests')
# test('basic', exe)

