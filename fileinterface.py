from io import TextIOWrapper
from pathlib import Path
BASE_DIR = Path(__file__).parent / 'resources'


def myopen(path, mode) -> TextIOWrapper:
    file = Path(path)
    file.parent.mkdir(parents=True, exist_ok=True)
    file.touch(exist_ok=True)
    return file.open(mode)


def html(text: str) -> str:
    with myopen(BASE_DIR / 'html.html', 'r') as f:
        template = f.read()
    return template.format(message=text)