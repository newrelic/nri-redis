# ###############################################################
# Settings and defaults.

# Load from flags or defaults from tilt_config
config.define_string("mode", args=True)
cfg = config.parse()

# Load the desired mode
symbols = load_dynamic("./tools/tilt-%s.starlark" % cfg["mode"])
run = symbols["run"]


# Setting up defaults and overrides for the selected mode
settings = {}
settings.update(read_json("tools/tilt-%s.json" % "global",    default={}))
settings.update(read_json("tools/tilt-%s.json" % cfg["mode"], default={}))
settings.update(read_json("tools/tilt-%s.json" % "local",     default={}))

default_registry(settings.get("default_registry"))

allow_k8s_contexts(settings.get("contexts", []))


# Run the desired mode
load("./tools/tilt-global.starlark", global_run="run")

global_run(settings)
run(settings)
