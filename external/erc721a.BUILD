load("@//sol:defs.bzl", "sol_sources")

package(default_visibility = ["//visibility:public"])

sol_sources(
    name = "erc721a",
    srcs = glob(["**/*.sol"]),
    remappings = {"erc721a/": "./external/erc721a/"},
)
