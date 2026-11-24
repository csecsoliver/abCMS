import sqlite3
import pathlib
import os
import json

db_conn = sqlite3.connect('abcms.sqlite')
cursor = db_conn.cursor()

# os.system('sh setup.sh')
users_dir = pathlib.Path('users')
for user_file in users_dir.iterdir():
    with open(user_file, 'r') as f:
        data = f.read()
    username = user_file.stem
    user = json.loads(data)
    social = user.get('social', "")
    upload_token = user.get('upload_token', 'None')
    if username in cursor.execute('SELECT username FROM users').fetchall():
        cursor.execute("UPDATE users SET xp = xp + ?, coins = coins + ?, social = ?, upload_token = ? WHERE username = ?",
                       (user.get('xp', user['coins']), user['coins'], social, upload_token, username))
        continue
    cursor.execute("INSERT INTO users (username, passhash, xp, coins, social, upload_token) VALUES (?, ?, ?, ?, ?, ?)",
                   (username, user['password'], user.get('xp', user['coins']), user['coins'], social, upload_token))
    db_conn.commit()



os.system("cp -r cozy static/uploads/")
with open('static/uploads/cozy/posts.json', 'r') as cozy_posts_file:
    cozy_posts_data = json.loads(cozy_posts_file.read())
counter = 0
for post in cozy_posts_data['posts']:
    cursor.execute("INSERT INTO posts (title, user_id, created_at, path, has_image) VALUES (?, (SELECT id FROM users WHERE username = ?), ?, ?, 1)",
                   (post['title'], post['path'].split('/')[1], counter, "/static/uploads/"+post['path']))
    counter += 1
    db_conn.commit()


markdown_dir = pathlib.Path('md_blog_content')
for post_file in markdown_dir.iterdir():
    if not post_file.suffix == '.md':
        continue
    with open(post_file, 'r') as f:
        content = f.readlines()
    post_id = post_file.stem
    with open(f"md_blog_content/{post_id}.json", 'r') as f:
        try:
            json_data = json.loads(f.read())
            color = json_data.get('color', '')
        except json.JSONDecodeError:
            color = ''
    title: str = content[0].replace("# ", "").strip()
    author: str = content[2].strip()
    content_md: str = "\n".join(content[4:])
    cursor.execute("INSERT INTO posts (title, content, user_id, created_at, color) VALUES (?, ?, (SELECT id FROM users WHERE username = ?), ?, ?)",
                   (title, content_md, author, counter + int(post_id), color))
    db_conn.commit()


db_conn.close()