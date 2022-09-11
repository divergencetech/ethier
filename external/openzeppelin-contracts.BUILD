load("@//sol:defs.bzl", "sol_sources")

package(default_visibility = ["//visibility:public"])

sol_sources(
    name = "openzeppelin-contracts",
    srcs = glob(["**/*.sol"]),
    remappings = {"@openzeppelin/contracts/": "./external/openzeppelin-contracts/contracts/"},
)
