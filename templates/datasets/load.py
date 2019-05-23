import pandas as pd
import os


def load():
    """Reads csv file if it exists, otherwise raises an exception"""
    this_dir = os.path.dirname(os.path.realpath(__file__))
    destfile = os.path.join(this_dir, "{{.RelPath}}")

    if (os.path.isfile(destfile)):
        data = pd.read_csv(destfile)
    else:
        raise Exception("Dataset missing from disk")

    return data
