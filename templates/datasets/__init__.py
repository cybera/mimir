import os

def list():
    """list all importable datasets"""
    path = os.path.dirname(os.path.abspath(__file__))

    for _, _, files in os.walk(path):
        for file in files:
            if file.endswith('.py') and file != '__init__.py':
                print(file[:-3]) 
