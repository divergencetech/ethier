"Implementation for sol_sources rule"
load("@aspect_rules_js//js:providers.bzl", "JsInfo", "js_info")
load("//sol:providers.bzl", "SolSourcesInfo")
load("@aspect_rules_js//js:libs.bzl", "js_lib_helpers")
load("@bazel_skylib//lib:dicts.bzl", "dicts")

_ATTRS = {
    "srcs": attr.label_list(
        allow_files = [".sol"],
        doc = "Solidity source files",
    ),
    "deps": attr.label_list(
        doc = "Each dependency should either be more .sol sources, or npm packages for 3p dependencies",
        providers = [[SolSourcesInfo], [JsInfo]],
    ),
    "remappings": attr.string_dict(
        doc = """Contribute to import mappings.
        
        See https://docs.soliditylang.org/en/latest/path-resolution.html?highlight=remappings#import-remapping
        """,
        default = {},
    ),
}

def _sol_sources_impl(ctx):
    npm_linked_packages = js_lib_helpers.gather_npm_linked_packages(
        srcs = ctx.attr.srcs,
        deps = ctx.attr.deps,
    )

    transitive_remappings = [dep[SolSourcesInfo].remappings for dep in ctx.attr.deps if SolSourcesInfo in dep]
    # TODO: detect and error on conflicting mappings from same value to different keys
    remappings = dicts.add(ctx.attr.remappings, *transitive_remappings)

    return [
        DefaultInfo(
            files = depset(ctx.files.srcs),
        ),
        SolSourcesInfo(
            direct_sources = ctx.files.srcs,
            transitive_sources = depset(
                ctx.files.srcs,
                transitive = [
                    d[SolSourcesInfo].transitive_sources
                    for d in ctx.attr.deps
                    if SolSourcesInfo in d
                ]
            ),
            remappings = remappings,
        ),
        js_info(
            npm_linked_packages = npm_linked_packages.direct,
            npm_linked_package_files = npm_linked_packages.direct_files,
            transitive_npm_linked_package_files = npm_linked_packages.transitive_files,
            transitive_npm_linked_packages = npm_linked_packages.transitive,
        ),
    ]

sol_sources = struct(
    implementation = _sol_sources_impl,
    attrs = _ATTRS,
)
