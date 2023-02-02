
cmake_minimum_required(VERSION 3.16)
project(rest LANGUAGES C VERSION 0.1)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
enable_testing()

add_library(mongoose inc/mongoose/mongoose.c)

add_executable(socks5-c main.c)

target_include_directories(webserver PRIVATE inc/mongoose)


# CAAA vector
target_link_libraries(webserver vector)
