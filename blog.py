import auth
import markdown
import random
from fileinterface import myopen
from pathlib import Path
def post(title: str, author: str, content: str, color:str,  existing_id: str = None) -> str:
    
    currentid = myopen('md_blog_content/currentid.txt', 'r').read()
    if currentid.isnumeric() == False:
        currentid = '0'
    newid = str(int(currentid)+1)
    if existing_id is not None:
        newid = existing_id
    content = f"# {title}\n\n{author}\n\n{color}\n\n{content}"
    
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
    color: str = content.split('\n')[4].strip()
    
    content_md: str = '\n'.join(content.split('\n')[6:])
    
    if format == 'html':
        return {'content': markdown.markdown(content_md), 'title': title, 'author':  f"{author}, who has {auth.get_coins(author)} coins", 'color': color}
    else:
        return {'content': content_md, 'title': title, 'author': author, 'color': color}
    
    
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

