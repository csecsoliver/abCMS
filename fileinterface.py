from pathlib import Path
import json
from typing import IO, Any

BASE_DIR = Path(__file__).parent / "resources"


def myopen(path: Path | str, mode: str) -> IO[str]:
    file = Path(path)
    file.parent.mkdir(parents=True, exist_ok=True)
    file.touch(exist_ok=True)
    return file.open(mode)


def getjson(path: Path | str) -> dict[str, Any]:
    try:
        with myopen(path, "r") as f:
            data = json.load(f)
            return data
    except json.JSONDecodeError:
        return {}


def setjson(path: Path | str, data: dict[str, str]):
    with myopen(path, "w") as f:
        json.dump(data, f)


def html(text: str) -> str:
    with myopen(BASE_DIR / "html.html", "r") as f:
        template = f.read()
    return template.format(message=text)
