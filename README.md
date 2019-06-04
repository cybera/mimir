# README

## Install

Download the latest version for your operating system from the [releases](https://github.com/cybera/ccds/releases) page and add it to your PATH.

## Usage

### Initialize

To initialize a new project in the current directory, run:

```bash
ccds init
```

The command will then ask a few questions about how to configure the project.

### Datasets

Create a new raw dataset by running:

```bash
ccds dataset new raw_dataset.csv
```

Note that file extension is currently required to detect the format, though only csv is currently supported.

Alternatively, to create a generated dataset that depends on a raw dataset, run:

```bash
ccds dataset new generated_dataset.csv -d="raw_dataset"
```

Note that the file extension is optional when declaring dependencies.

### Jupyter

The `jupyter` subcommand manages the Jupyter Labs Docker container:

```bash
ccds jupyter <start/stop>
```

### Scripts

Any scripts in the `src/scripts` directory can be run with the command:

```bash
ccds run <script_name>
```

Scripts in that directory should be executable and have a correct shebang for their filetype (e.g. `#!/usr/bin/env python`)

## Settings

When initializing a new project, a `project-settings.toml.example` is created in the project root. It is intended to be a template for creating project-wide settings that can be used in your code. If you want to use it, modify it as needed with any settings and default values and then commit it to the repository. It should be copied to `project-settings.toml` when the repository is cloned and any remaining values filled in. The resulting file should never be committed in order to avoid leaking secrets.

How to access the settings depends on the project's language.

For Python:

```python
from src import settings

settings.settings["setting_name"]
```

For R:

```R
source("path/to/src/settings.R")

settings["setting_name"]
```
