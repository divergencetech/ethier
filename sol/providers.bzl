"Providers for rule interop"

SolSourcesInfo = provider(
    doc = "Stores a tree of source file dependencies",
    fields = {
        "direct_sources": "list of sources provided to this node",
        "transitive_sources": "depset of transitive dependency sources",
        "remappings": "dictionary of import remappings to propagate to solc",
    },
)

