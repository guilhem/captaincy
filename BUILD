load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "controller.go",
        "main.go",
    ],
    importpath = "github.com/guilhem/captaincy",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/api/apps/v1beta2:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/client-go/informers:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/typed/core/v1:go_default_library",
        "//vendor/k8s.io/client-go/listers/apps/v1beta2:go_default_library",
        "//vendor/k8s.io/client-go/tools/cache:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/k8s.io/client-go/tools/record:go_default_library",
        "//vendor/k8s.io/client-go/util/workqueue:go_default_library",
        "//vendor/github.com/guilhem/captaincy/pkg/apis/samplecontroller/v1alpha1:go_default_library",
        "//vendor/github.com/guilhem/captaincy/pkg/client/clientset/versioned:go_default_library",
        "//vendor/github.com/guilhem/captaincy/pkg/client/clientset/versioned/scheme:go_default_library",
        "//vendor/github.com/guilhem/captaincy/pkg/client/informers/externalversions:go_default_library",
        "//vendor/github.com/guilhem/captaincy/pkg/client/listers/samplecontroller/v1alpha1:go_default_library",
        "//vendor/github.com/guilhem/captaincy/pkg/signals:go_default_library",
    ],
)

go_binary(
    name = "sample-controller",
    importpath = "github.com/guilhem/captaincy",
    library = ":go_default_library",
    visibility = ["//visibility:public"],
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
        "//staging/src/github.com/guilhem/captaincy/pkg/apis/samplecontroller:all-srcs",
        "//staging/src/github.com/guilhem/captaincy/pkg/client/clientset/versioned:all-srcs",
        "//staging/src/github.com/guilhem/captaincy/pkg/client/informers/externalversions:all-srcs",
        "//staging/src/github.com/guilhem/captaincy/pkg/client/listers/samplecontroller/v1alpha1:all-srcs",
        "//staging/src/github.com/guilhem/captaincy/pkg/signals:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
