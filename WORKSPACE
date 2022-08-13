load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")

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

new_git_repository(
    name = "com_github_rivo_uniseg",
    build_file_content = """load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "uniseg",
    srcs = [
        "doc.go",
        "eastasianwidth.go",
        "grapheme.go",
        "graphemeproperties.go",
        "graphemerules.go",
        "line.go",
        "lineproperties.go",
        "linerules.go",
        "properties.go",
        "sentence.go",
        "sentenceproperties.go",
        "sentencerules.go",
        "step.go",
        "word.go",
        "wordproperties.go",
        "wordrules.go",
    ],
    importpath = "github.com/rivo/uniseg",
    visibility = ["//visibility:public"],
)
""",
    commit = "56f4d68f787bd41bbd705a494c955d48784ace5f",  # v0.3.4
    remote = "https://github.com/rivo/uniseg.git",
    shallow_since = "1659733471 +0200",
)

# Tink's Go implementation is one directory down and git_repository doesn't have
# an obvious way to traverse, so we use strip_prefix instead.
TINK_VERSION = "1.7.0"

TINK_URL = "https://github.com/google/tink/archive/refs/tags/v{}.zip".format(TINK_VERSION)

TINK_SHA256 = "ff272c968827ce06b262767934dc56ab520caa357a4747fc4a885b1cc711222f"

http_archive(
    name = "tink_base",
    sha256 = TINK_SHA256,
    strip_prefix = "tink-{}".format(TINK_VERSION),
    url = TINK_URL,
)

http_archive(
    name = "com_github_google_tink_go",
    sha256 = TINK_SHA256,
    strip_prefix = "tink-{}/go".format(TINK_VERSION),
    url = TINK_URL,
)

#############
# Go Ethereum
#############

GO_ETHEREUM_SRC = "github.com/ethereum/go-ethereum"

GO_ETHEREUM_VERSION = "1.10.21"

http_archive(
    name = "geth_secp256k1",
    build_file_content = """
cc_library(
    name = "hdrs",
    hdrs = glob([
        "**/*.c",
        "**/*.h",
    ]),
    visibility = ["//visibility:public"],
)
""",
    sha256 = "ad8ffdce7dd530ada33bffcd6d4a1005c57e404ef8a56730f4af73cb74c0d97d",
    strip_prefix = "go-ethereum-{}/crypto/secp256k1".format(GO_ETHEREUM_VERSION),
    url = "https://{}/archive/refs/tags/v{}.zip".format(GO_ETHEREUM_SRC, GO_ETHEREUM_VERSION),
)

go_repository(
    name = "com_github_ethereum_go_ethereum",
    importpath = "github.com/ethereum/go-ethereum",
    patches = ["//patches:com_github_ethereum_go_ethereum.patch"],  # keep
    sum = "h1:5lqsEx92ZaZzRyOqBEXux4/UR06m296RGzN3ol3teJY=",
    version = "v1.10.21",
)

# From go.mod; update with:
# $ bazel run //:gazelle update-repos -- -from_file=go.mod -to_macro=repositories.bzl%go_repositories
load("//:repositories.bzl", "go_repositories")

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

go_rules_dependencies()

go_register_toolchains(version = "1.18")

gazelle_dependencies()
