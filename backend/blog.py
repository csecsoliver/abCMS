import markdown
import random
from fileinterface import myopen
from pathlib import Path
def post(title: str, content):
    currentid = myopen('md_blog_content/currentid.txt', 'r').read()
    if currentid.isnumeric() == False:
        currentid = '0'
    newid = str(int(currentid)+1)
    myopen(f'md_blog_content/{newid}.md', 'w').write(content)
    
    myopen('md_blog_content/currentid.txt', 'w').write(newid)
    return newid
    
    
def get(id: str, format: str='md'):
    if id.isnumeric() == False:
        return 'Invalid ID'
    content = myopen('md_blog_content/'+id+'.md', 'r').read()
    if format == 'html':
        return markdown.markdown(content)
    else:
        return content
    
    
def listids(reverse: bool=False):
    ids = []
    folder = Path('md_blog_content')
    files = folder.iterdir()
    for file in files:
        if file.suffix == '.md' and file.stem.isnumeric():
            ids.append(file.stem)
    ids.sort(key=lambda x: int(x), reverse=reverse)
    return ids