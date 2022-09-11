"""# Bazel rules for Solidity

See <https://docs.soliditylang.org>
"""

load("//sol/private:sol_binary.bzl", lib = "sol_binary")
load("//sol/private:sol_sources.bzl", src = "sol_sources")
load(":providers.bzl", "SolSourcesInfo")

sol_binary = rule(
    implementation = lib.implementation,
    attrs = lib.attrs,
    doc = """sol_binary compiles Solidity source files with solc""",
    toolchains = lib.toolchains,
)

sol_sources = rule(
    implementation = src.implementation,
    attrs = src.attrs,
    doc = """Collect .sol source files to be imported as library code.
    Performs no actions, so semantically equivalent to filegroup().
    """,
    provides = [SolSourcesInfo],
)
