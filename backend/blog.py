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
    
    
def get(id: str, format: str='md') -> dict:
    if id.isnumeric() == False:
        return 'Invalid ID'
    fileob = myopen('md_blog_content/'+id+'.md', 'r')
    content = fileob.read()
    fileob.close()
    title: str = content.split('\n')[0].replace('#', '').strip()
    author: str = content.split('\n')[2].strip()
    content_md: str = '\n'.join(content.split('\n')[4:])
    
    if format == 'html':
        return {'content': markdown.markdown(content_md), 'title': title, 'author': author}
    else:
        return {'content': content_md, 'title': title, 'author': author}
    yx
    
def listids(reverse: bool=False):
    ids = []
    folder = Path('md_blog_content')
    if not folder.exists():
        return ids
    files = folder.iterdir()
    for file in files:
        if file.suffix == '.md' and file.stem.isnumeric():
            ids.append(file.stem)
    ids.sort(key=lambda x: int(x), reverse=reverse)
    return ids