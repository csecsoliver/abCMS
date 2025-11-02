"""
Contains methods related to managin the text based blog posts.
"""
import json
import auth
import markdown
from fileinterface import myopen
from pathlib import Path
import bleach
import webcolors
import colorsys


def post(
    title: str, author: str, content: str, color: str, existing_id: str = ""
) -> str:
    """Create post with the specified content

    Args:
        title (str): Title string to set. Does NOT get sanitized
        author (str): Author string to set. Does NOT get sanitized
        content (str): Content string to set. Does NOT get sanitized
        color (str):  Color string to set. Does NOT get sanitized or validated
        existing_id (str, optional): Existing post to modify. Not tested. Defaults to "".

    Returns:
        str: The id of the new post
    """
    with myopen("md_blog_content/currentid.txt", "r") as f:
        currentid = f.read()
    if currentid.isnumeric() == False:
        currentid = "0"
    newid = str(int(currentid) + 1)
    if existing_id != "":
        newid = existing_id
    content = f"# {title}\n\n{author}\n\n{content}"

    with myopen(f"md_blog_content/{newid}.md", "w") as f:
        _ = f.write(content)
    with myopen(f"md_blog_content/{newid}.json", "w") as f:
        json.dump({"color": color}, f)
    with myopen("md_blog_content/currentid.txt", "w") as f:
        _ = f.write(newid)
    return newid


def get(id: str, format: str = "md", truncate: bool = False) -> dict[str, str]:
    """Gets the content of a blog post

    Args:
        id (str): The post id to look for
        format (str, optional): The format of the returned string. Can be "html" or "md" Defaults to "md".
        truncate (bool, optional): Whether to truncate the body of the post at 200 characters. Defaults to False.

    Returns:
        dict[str, str]: Post data. Contains content, title, author and color.
    """
    if id.isnumeric() == False:
        return {
            "content": "Post not found.",
            "title": "Interesting...",
            "author": "",
            "color": "#ffffff",
        }

    with myopen("md_blog_content/" + id + ".md", "r") as fileob:
        content = fileob.readlines()

    with myopen(f"md_blog_content/{id}.json", "r") as f:
        try:
            meta = json.load(f)
        except:
            meta = {}

    title: str = content[0].replace("# ", "").strip()
    author: str = content[2].strip()

    content_md: str = "\n".join(content[4:])
    if truncate and len(content_md) > 200:
        content_md = content_md[:200] + "...\n\n*Click the title to read more.*"

    if format == "html":
        content_md = bleach.clean(content_md)
        content_md = markdown.markdown(
            content_md, extensions=["extra", "nl2br", "sane_lists", "wikilinks"]
        )
        content_md = bleach.linkify(content_md, skip_tags=["pre"])
        title = bleach.clean(title)
        return {
            "content": content_md,
            "title": title,
            "author": f"{author}<br>MANA={auth.get_coins(author)}, XP={auth.get_xp(author)}",
            "color": meta.get("color", "#ffffff"),
        }
    else:
        return {
            "content": content_md,
            "title": title,
            "author": author,
            "color": meta.get("color", "#ffffff"),
        }


def listids(reverse: bool = False, username: str = "") -> list[str]:
    """List the existing posts by id

    Args:
        reverse (bool, optional): Whether to reverse the list. Normal order is oldest first. Defaults to False.
        username (str, optional): Filter posts by username. If it's empty, they are not filtered. Defaults to "".

    Returns:
        list[str]: List of id strings (actually always numbers)
    """
    ids = []
    folder = Path("md_blog_content")
    if not folder.exists():
        return ids
    files = folder.iterdir()
    for file in files:
        if file.suffix == ".md" and file.stem.isnumeric():
            ids.append(file.stem)
    ids.sort(key=lambda x: int(x), reverse=reverse)
    if username != "":
        filtered_ids = []
        for id in ids:
            post_dict = get(id, format="md")
            if post_dict["author"] == username:
                filtered_ids.append(id)
        return filtered_ids
    return ids


def delete(id: str, username: str) -> bool:
    """Delete Post by id and username

    Args:
        id (str): Id of the post
        username (str): Username of user requesting the action

    Returns:
        bool: Whether is succeeded
    """
    post_dict: dict[str, str] = get(id, format="md")
    if post_dict["author"] != username:
        return False
    Path(f"md_blog_content/{id}.md").unlink(missing_ok=True)
    return True


def get_text_color(bg_color: str):
    """Handles the setting of post colors. I am not sure this works correctly

    Args:
        bg_color (str): CSS compatible color representation of the background

    Returns:
        str: CSS compatible representation of the best text color for it
    """
    htmlcolor = webcolors.html5_parse_legacy_color(bg_color)
    lightness = colorsys.rgb_to_hls(*htmlcolor)[1]  # HLS lightness (0-1)
    print(lightness)
    return "#000" if lightness > 130 else "#ccc"
