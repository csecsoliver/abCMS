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


def html(text: str, back="/", cozy=False) -> str:
    with myopen(BASE_DIR / ("html.html" if not cozy else "cozyhtml.html"), "r") as f:
        template = f.read()
    return template.format(message=text, back=back)

def clean(s:str)->str:
    s = "".join(c for c in s if c.isalnum() or c in ("_", "-", ".")).strip()
    return s


