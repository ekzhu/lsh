add_custom_target(project_test ALL)

add_executable(main main.cpp)
target_link_libraries(main flann_cpp)
add_dependencies(project_test main)

