# CCDS

CCDS is an interactive CLI tool that simplifies common repetitive tasks and provides a logical and consistent structure for your data science projects. 

Pre-configured for use with Git, Docker, and Jupyter Labs. Currently supports R and Python development.

### Contents
1. [Download & Installation](#install)
2. [Creating A New Project](#create)
3. [Configuring Datasets](#adding)
4. [Fetching & Updating Datasets](#fetching)
5. [Jupyter Notebooks](#jupyter)
6. [Custom Scripts](#scripts)
6. [Project Settings](#settings)

&nbsp;

<a name="install" />

## Download & Installation

### Dependencies
If you wish to use the Jupyter Labs feature, ensure that you have [Docker](https://docs.docker.com/install/) installed on your computer.

### Latest Release
Download it [here](https://github.com/cybera/ccds/releases).

### Development Version

1. `git clone git@github.com:cybera/ccds.git` and checkout the branch you wish to use. 
2. Ensure you've installed [golang](https://golang.org/dl/) and set the [$GOPATH variable](https://github.com/golang/go/wiki/SettingGOPATH#unix-systems) (it will be set automatically if you're installing on Windows)
4. Install packr2 by running `go get -u github.com/gobuffalo/packr/v2/packr2`. 
5. Run `packr2 install` from the CCDS root directory.

&nbsp;

<a name="create" />

## Creating A New Project
CCDS makes it simple to move between different projects since each project will follow the same structure. Regardless of what you're working on or who you're collaborating with, you'll always know exactly where to find everything.

Running the `ccds init` command from a new directory and following the project set-up workflow will create the following directories:

```bash
├───.ccds
├───data
│   ├───external
│   ├───interim
│   ├───processed
│   └───raw
├───docs
├───models
├───notebooks
├───references
├───reports
│   └───figures
└───src
    ├───datasets
    ├───features
    ├───models
    ├───scripts
    └───visualization
```

Other files created include a `LICENSE` file, an initialized git repository, and initial Jupyter configuration for Docker to run.

&nbsp;

<a name="adding" />

## Configuring Datasets
Before you can jump into working with your datasets, you need to let CCDS know where it can find them.

Please note, CSV is the only file format currently supported for datasets. You can view and edit dataset configuration at `.ccds/project-metadata.toml`.

### Stored on your computer
If you have a raw dataset saved on your computer that you'd like to add, run the following command:

```bash
ccds dataset add my_dataset.csv --from ~/Downloads/source_data.csv
```

This will make CCDS aware of the dataset's source location and generate import code in `src/datasets/my_dataset.py` or `src/datasets/my_dataset.r` so that you can start working with your dataset right away. This will not copy the dataset yet, just creating a R/Python template for it.

### Stored somewhere else
You can import datasets stored in an external cloudprovider using OpenStack. You'll need to ensure that you've saved your OpenStack security credentials and sourced them. Currently CCDS only supports fetching through Swift, so [ensure you have it installed](https://swift.org/getting-started/#using-the-repl)

Let's say I have a container called `datasets_container` that hosts my source dataset, `cloud_source.csv`.

```bash
ccds dataset add cloud_data.csv --source swift --from datasets_container/cloud_source.csv
```

### Generated from an existing dataset
If you'd like to create a new dataset that depends on existing raw data that you've already added, run the following command:

```bash
ccds dataset add generated_dataset.csv --generated --depends-on=my_dataset
```

A new file will be created at `/src/datasets/generated_dataset.py` or `/src/datasets/generated_dataset.r`.

&nbsp;

<a name="fetching" />

## Fetching & Updating Datasets

Fetching datasets requires adding them first via `ccds dataset add` so that CCDS knows where they're being pulled from.

### All datasets at once
After adding all of the sources for your datasets via `ccds dataset add`, you can easily download all of them into their appropriate folders at once by running the following command:

```bash
ccds dataset fetch
```

The CLI will ask to confirm that all datasets should be fetched. Enter `y` to continue. All datasets will be copied and sorted into the following `/data/` subdirectories:

```bash
├───data
│   ├───processed # Generated datasets
│   └───raw       # Local & external datasets
```

### Specific dataset
Perhaps you want to download one specific dataset. In that case, you just need to specify the name of the added dataset without `.csv`.

```bash
ccds dataset fetch my_data
```

Since I added `my_dataset` from a local source, I would be able to find the resulting dataset at `/data/raw/my_dataset.csv`.

### Updating a fetched dataset
Let's say your source dataset has been modified. You can bring your CCDS project dataset up to date by refetching it and overwriting the old version.

```bash
ccds dataset fetch my_dataset
```

The CLI will respond with `my_data already exists, fetch a new copy? [y/N]`. Enter `y` to continue.

### Troubleshooting

#### Error: The system cannot find the specified file
Double-check that your file exists at the source you specified during the `ccds dataset add` command step. If everything's fine and dandy there, try also checking the CCDS configuration file at `.ccds/project-metadata.toml`.

For example, let's say we made a typo in the file name, entering `src_data.csv` instead of `source_data.csv`.

```bash
ccds dataset add my_dataset.csv --from ~/Downloads/src_data.csv
```

This is what you would see in the metadata:

```toml
[datasets]
  [datasets.my_dataset]
    dependencies = []
    file = "my_dataset.csv"
    generated = false
    [datasets.my_dataset.source]
      name = "local"
      target = "Downloads\\src_data.csv" <-- The line to edit
```

Under `[datasets.my_dataset.source]`, edit the `target` line. Rerun `ccds dataset fetch` to try pulling the dataset again. If everything is good to go, you should now see your dataset appear in the correct `/data/` subdirectory.

&nbsp;

<a name="jupyter" />

## Jupyter Notebooks
The `jupyter` subcommand manages the Jupyter Labs Docker container. 

### Starting Jupyter
To start a Jupyter Labs instance, run the following command:

```bash
ccds jupyter start
```

Once the Docker container for Jupyter is up and running, you will be able to access it by following one of the links provided in the CLI. One will be for the notebooks app and the other will be for Jupyter Lab.

The Jupyter instance manages all files in the `notebooks` project directory. You can add your own `.ipynb` files to this directory in order to access and edit them.

### Quitting Jupyter
Stop the container at any time by running the stop command:

```bash
ccds jupyter stop
```

### Troubleshooting

#### Error: no connection, localhost:8888 won't load in the browser 
You may need to [open your hosts configuration](https://en.wikiversity.org/wiki/Hosts_file/Edit) and double-check that the following line exists:

```bash
127.0.0.1   localhost
```

&nbsp;

<a name="scripts" />

## Custom Scripts
You can run your own executable custom bash, python, and R scripts with CCDS.

### Set-up
1. Save your file to the `src/scripts/` directory.
2. Add a [she-bang line](https://linoxide.com/linux-how-to/what-function-shebang-linux/) to the top of the file if you haven't already. For example, `#!/bin/bash` for bash scripts or `#!/usr/bin/env python` for python.
3. From your terminal, run `cd src/scripts && chmod u+x my_script_file.py`

### Running A Script
Let's say you followed the set-up instructions and named my script `my_script_file.py`. To run it, you would enter the following command:

```bash
ccds run my_script_file.py
```

&nbsp;

<a name="settings" />

## Project Settings
All of your project's settings can be managed by editing `project-settings.toml`.

### Custom Settings Template
Running the `ccds init` command will create a `project-settings.toml.example` file in the root directory of your project. It is intended to be a template for creating project-wide settings that can be used in your code. If you want to use it, add any settings and default values and then commit it to the repository.

Settings saved to `/project-settings.toml.example` will have to be manually copied to `/.ccds/project-settings.toml` after the repository is cloned.

### Using Settings in Code
You can use settings to set default values for variables in your project code. Importing the settings depends on the programming language you're using. They can be used in scripts and Jupyter notebooks the same way.

**For Python:**

```python
from src import settings

settings.settings["setting_name"]
```

**For R:**

```r
source("../src/settings.R") # You may need to edit the path to the source

settings["setting_name"]
```
