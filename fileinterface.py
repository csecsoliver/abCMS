from io import TextIOWrapper
from pathlib import Path
def myopen(path, mode) -> TextIOWrapper:
    file = Path(path)
    file.parent.mkdir(parents=True, exist_ok=True)
    file.touch(exist_ok=True)
    return file.open(mode)