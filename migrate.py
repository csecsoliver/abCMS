import sqlite3
import pathlib
import os
import json

db = sqlite3.connect('abcms.sqlite')
cursor = db.cursor()

# os.system('sh setup.sh')
usersdir = pathlib.Path('users')
users = usersdir.iterdir()
for username in users:
    with open(username, 'r') as f:
        data = f.read()
    username = username.stem
    user = json.loads(data)
    social = user.get('social', "")
    upload_token = user.get('upload_token', 'None')
    if username in cursor.execute('SELECT username FROM users').fetchall():
        cursor.execute("UPDATE users SET xp = xp + ?, coins = coins + ?, social = ?, upload_token = ? WHERE username = ?",
                       (user.get('xp', user['coins']), user['coins'], social, upload_token, username))
        continue
    cursor.execute("INSERT INTO users (username, passhash, xp, coins, social, upload_token) VALUES (?, ?, ?, ?, ?, ?)",
                   (username, user['password'], user.get('xp', user['coins']), user['coins'], social, upload_token))
    db.commit()



os.system("cp -r cozy static/uploads/")
cozyposts = open('static/uploads/cozy/posts.json', 'r')

cozyposts_data = json.loads(cozyposts.read())
counter = 0
for post in cozyposts_data['posts']:
    cursor.execute("INSERT INTO posts (title, user_id, created_at, path, has_image) VALUES (?, (SELECT id FROM users WHERE username = ?), ?, ?, 1)",
                   (post['title'], post['path'].split('/')[1], counter, "/static/uploads/"+post['path']))
    counter += 1
    db.commit()


mddir = pathlib.Path('md_blog_content')
posts = mddir.iterdir()
for postidmd in posts:
    if not postidmd.suffix == '.md':
        continue
    with open(postidmd, 'r') as f:
        content = f.readlines()
    postidmd = postidmd.stem
    with open(f"md_blog_content/{postidmd}.json", 'r') as f:
        try:
            json_data = json.loads(f.read())
            color = json_data.get('color', '')
        except json.JSONDecodeError:
            color = ''
    title: str = content[0].replace("# ", "").strip()
    author: str = content[2].strip()
    content_md: str = "\n".join(content[4:])
    cursor.execute("INSERT INTO posts (title, content, user_id, created_at, color) VALUES (?, ?, (SELECT id FROM users WHERE username = ?), ?, ?)",
                   (title, content_md, author, counter + int(postidmd), color))
    db.commit()


db.close()