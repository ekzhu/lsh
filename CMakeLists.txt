cmake_minimum_required(VERSION 2.6)

project(CSC2515Project)
string(TOLOWER ${PROJECT_NAME} PROJECT_NAME_LOWER)

# detect if using the Clang compiler
if("${CMAKE_C_COMPILER_ID}" MATCHES "Clang")
  set(CMAKE_COMPILER_IS_CLANG 1)
endif ()

if("${CMAKE_CXX_COMPILER_ID}" MATCHES "Clang")
  set(CMAKE_COMPILER_IS_CLANGXX 1)
endif ()

#set the default path for built executables to the "bin" directory
set(EXECUTABLE_OUTPUT_PATH "${PROJECT_BINARY_DIR}/bin")
#set the default path for built libraries to the "lib" directory
set(LIBRARY_OUTPUT_PATH "${PROJECT_BINARY_DIR}/lib")

# Include Flann
include("${PROJECT_SOURCE_DIR}/flann/cmake/FindFlann.cmake")

#set the C/C++ include path to the "include" directory
include_directories(BEFORE "${PROJECT_SOURCE_DIR}/flann/src/cpp")

# require proper c++
#add_definitions( "-Wall -ansi -pedantic" )
# HDF5 uses long long which is not ansi
if(CMAKE_C_COMPILER_ID MATCHES "MSVC" OR CMAKE_CXX_COMPILER_ID MATCHES "MSVC")
    # lots of warnings with cl.exe right now, use /W1
    add_definitions("/W1 -D_CRT_SECURE_NO_WARNINGS -D_SCL_SECURE_NO_WARNINGS /bigobj")
else()
    add_definitions( "-Wall -Wno-unknown-pragmas -Wno-unused-function" )
endif()

add_subdirectory(src)

