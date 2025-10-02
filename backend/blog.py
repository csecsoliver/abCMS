import markdown
import random
def post(title: str, content):
    currentid = open('md_blog_content/currentid.txt', 'r').read()
    newid = str(int(currentid)+random.randint(50,100))
    open(f'md_blog_content/{newid}.md', 'w').write(content)
    
    open('md_blog_content/currentid.txt', 'w').write(newid)
    return newid
    
def get(id: str, format):
    if id.isnumeric() == False:
        return 'Invalid ID'
    content = open('md_blog_content/'+id+'.md', 'r').read()
    if format == 'html':
        return markdown.markdown(content)
    return content
