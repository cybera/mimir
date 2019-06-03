import os
import toml

_this_dir = os.path.dirname(os.path.realpath(__file__))
_settings_file = os.path.join(_this_dir, "{{.ProjectSettingsPath}}")

settings = toml.load(_settings_file)
