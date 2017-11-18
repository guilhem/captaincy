load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["interface.go"],
    importpath = "github.com/guilhem/captaincy/pkg/client/informers/externalversions/samplecontroller",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/guilhem/captaincy/pkg/client/informers/externalversions/internalinterfaces:go_default_library",
        "//vendor/github.com/guilhem/captaincy/pkg/client/informers/externalversions/samplecontroller/v1alpha1:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//staging/src/github.com/guilhem/captaincy/pkg/client/informers/externalversions/samplecontroller/v1alpha1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
