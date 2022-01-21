# ###############################################################
# Settings and defaults.

settings = {}
settings.update(read_json("tilt_settings.json", default={}))
settings.update(read_json("devel/tilt_settings.json", default={}))

default_registry(settings.get("default_registry"))
allow_k8s_contexts(
    settings.get("allowed_contexts", settings.get("kind_cluster_name")),
)

# ###############################################################
# Reload Mechanism

load("ext://restart_process", "docker_build_with_restart")
docker_build_with_restart(
    settings.get("image_name", "%s-dev" % settings.get("project_name")),
    ".",
    dockerfile="devel/Dockerfile",
    entrypoint=settings.get("entrypoint"),
    live_update=[
        sync("./bin/%s" % settings.get("binary_name", settings.get("project_name")), "/var/db/newrelic-infra/newrelic-integrations/bin/%s" % settings.get("project_name")),
    ],
)

# ###############################################################
# Resources to create
# integration binary
local_resource(
    "integration-binary",
    "GOOS=linux make compile",
    deps=[
        "./src",
    ]
)

# integration chart deployment
k8s_yaml(
    helm(
        settings.get("chart_path", "devel/charts/integration"),
        name="integration-deployment",
        values=[
            "devel/values-local.yaml"
        ]
    )
)

k8s_resource(
    "integration-deployment",
    port_forwards = settings.get("port_forwards")
)
