# this file includes all day*.py files to run the asserts in them

import os
import importlib

current_dir = os.path.dirname(os.path.abspath(__file__))
files = [
    f[:-3] for f in os.listdir(current_dir) if f.startswith("day") and f.endswith(".py")
]

for file in files:
    importlib.import_module(file)
