import json
import auth
import markdown
import random
from fileinterface import myopen
from pathlib import Path
import bleach
def post(title: str, author: str, content: str, color:str,  existing_id: str = None) -> str:
    with myopen('md_blog_content/currentid.txt', 'r') as f:
        currentid = f.read()
    if currentid.isnumeric() == False:
        currentid = '0'
    newid = str(int(currentid)+1)
    if existing_id is not None:
        newid = existing_id
    content = f"# {title}\n\n{author}\n\n{content}"
    
    with myopen(f'md_blog_content/{newid}.md', 'w') as f:
        f.write(content)
    with myopen(f'md_blog_content/{newid}.json', 'w') as f:
        json.dump({'color': color}, f)
    with myopen('md_blog_content/currentid.txt', 'w') as f:
        f.write(newid)
    return newid
    
    
def get(id: str, format: str='md', truncate=False) -> dict:
    if id.isnumeric() == False:
        return 'Invalid ID'
    fileob = myopen('md_blog_content/'+id+'.md', 'r')

    with myopen(f'md_blog_content/{id}.json', 'r') as f:
        try:
            meta = json.load(f)
        except:
            meta = {}

    content = fileob.readlines()
    fileob.close()
    title: str = content[0].replace('# ', '').strip()
    author: str = content[2].strip()
    
    content_md: str = '\n'.join(content[4:])
    if truncate and len(content_md) > 200:
        content_md = content_md[:200] + "...\n\n*Click to read more.*"
    
    if format == 'html':
        
        content_md = bleach.clean(content_md)
        content_md = markdown.markdown(content_md, extensions=['extra', 'nl2br', 'sane_lists', 'wikilinks'])
        content_md = bleach.linkify(content_md, skip_tags=['pre'])
        return {'content': content_md, 'title': title, 'author':  f"{author}, COINS={auth.get_coins(author)}, XP={auth.get_xp(author)}", 'color': meta.get('color', '#ffffff')}
    else:
        return {'content': content_md, 'title': title, 'author': author, 'color': meta.get('color', '#ffffff')}
    
    
def listids(reverse: bool=False, username: str=None) -> list[str]:
    ids = []
    folder = Path('md_blog_content')
    if not folder.exists():
        return ids
    files = folder.iterdir()
    for file in files:
        if file.suffix == '.md' and file.stem.isnumeric():
            ids.append(file.stem)
    ids.sort(key=lambda x: int(x), reverse=reverse)
    if username is not None:
        filtered_ids = []
        for id in ids:
            post_dict = get(id, format='md')
            if post_dict['author'] == username:
                filtered_ids.append(id)
        return filtered_ids
    return ids

def delete(id: str, username: str) -> bool:
    post_dict = get(id, format='md')
    if post_dict['author'] != username:
        return False
    Path(f'md_blog_content/{id}.md').unlink(missing_ok=True)
    return True

