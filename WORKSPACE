load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")
load("//sol:repositories.bzl", "LATEST_VERSION", "rules_sol_dependencies", "sol_register_toolchains")

rules_sol_dependencies()

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

##########
# Protobuf
##########

http_archive(
    name = "com_google_protobuf",
    sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
    strip_prefix = "protobuf-3.14.0",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

########
# Golang
########

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "ab21448cef298740765f33a7f5acee0607203e4ea321219f2a4c85a6e0fb0a27",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.32.0/rules_go-v0.32.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.32.0/rules_go-v0.32.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "5982e5463f171da99e3bdaeff8c0f48283a7a5f396ec5282910b9e8a49c0dd7e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

#########
# Go deps
#########

# From go.mod; update with:
# `bazel run //:gazelle update-repos -- -from_file=go.mod -to_macro=repositories.bzl%go_repositories`
load("//:repositories.bzl", "go_repositories")

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

# Necessary for `//:gazelle update-repos`
go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:6zppjxzCulZykYSLyVDYbneBfbaBIQPYMevg0bEwv2s=",
    version = "v0.6.0-dev.0.20220419223038-86c51ed26bb4",
)

go_repository(
    name = "org_golang_x_xerrors",
    importpath = "golang.org/x/xerrors",
    sum = "h1:5Pf6pFKu98ODmgnpvkJ3kFUOQGGLIzLIkbzUHp47618=",
    version = "v0.0.0-20220517211312-f3a8303e98df",
)

# Unlike the other Go packages, these already have BUILD.bazel files so we
# don't need to use the go_repository() rule, which creates them.
git_repository(
    name = "com_github_h_fam_errdiff",
    commit = "784941bedd9a6e28369be52ef6c2d234f77c10fc",  # v1.0.2
    remote = "https://github.com/h-fam/errdiff.git",
    shallow_since = "1585675534 -0700",
)

git_repository(
    name = "com_github_bazelbuild_tools_jvm_autodeps",
    commit = "62694dd50b91955fdfe67d0c0583fd25d78e5389",
    remote = "https://github.com/bazelbuild/tools_jvm_autodeps.git",
    shallow_since = "1537169762 +0200",
)

go_rules_dependencies()

go_register_toolchains(version = "1.18")

gazelle_dependencies()



sol_register_toolchains(
    name = "solc",
    sol_version = LATEST_VERSION,
)

# Fetch our contracts dependencies directly from git repos, like Forge does it:
# https://book.getfoundry.sh/projects/dependencies
new_git_repository(
    name = "openzeppelin-contracts",
    remote = "git@github.com:OpenZeppelin/openzeppelin-contracts.git",
    # NOTE: this version ought to match what appears in yarn.lock, or behavior
    # between Bazel and legacy npm-based dependency management will vary.
    commit = "8c49ad74eae76ee389d038780d407cf90b4ae1de", # v4.7.0
    shallow_since = "1656493217 +0200",
    build_file = "openzeppelin-contracts.BUILD",
)

new_git_repository(
    name = "erc721a",
    remote = "git@github.com:chiru-labs/ERC721A.git",
    # NOTE: this version ought to match what appears in yarn.lock, or behavior
    # between Bazel and legacy npm-based dependency management will vary.
    commit = "9859cd2edb1a8b4c2db5e46031abbb3253e42467", # v4.2.2
    shallow_since = "1659313679 -0700",
    build_file = "erc721a.BUILD",
)
