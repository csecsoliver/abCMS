import argon2
import uuid
from bottle import Response
from typing import *

ph = argon2.PasswordHasher()
tokens = {}
def checkpass(username, password, response: Response) -> bool:
    with open(f"users/{username}", "r") as f:
        hashed = f.read()
    try:
        ph.verify(hashed, password)
        response.set_cookie("user", username)
        token = uuid.uuid4().hex
        tokens[token] = username
        response.set_cookie("token", token, secret=getsecret()["cs"], httponly=True, samesite='Strict')
        return True
    except argon2.exceptions.VerifyMismatchError:
        response.status = 401
        return False
    
    

def createuser(username, password, secret, response: Response) -> bool:
    if secret != getsecret()["ss"]:
        return False
    if open(f"users/{username}", "r"):
        return False
    hashed = ph.hash(password)
    with open(f"users/{username}", "w") as f:
        f.write(hashed)
    response.set_cookie("user", username)
    token = uuid.uuid4().hex
    tokens[token] = username
    response.set_cookie("token", token, secret=getsecret()["cs"], httponly=True, samesite='Strict')
    return True

def getsecret(opt: Literal["ss", "cs"]) -> str:
    signupsecret = "secrete"
    cookiesecret = "asdfkjsdd"
    with open(".env", "r") as f:
        for line in f:
            if line.startswith("SIGNUP_SECRET="):
                signupsecret = line[len("SIGNUP_SECRET="):].strip()
            elif line.startswith("COOKIE_SECRET="):
                cookiesecret = line[len("COOKIE_SECRET="):].strip()
    if signupsecret == "secrete":
        print("SIGNUP_SECRET not found in .env, using insecure default.")
    if cookiesecret == "asdfkjsdd":
        print("COOKIE_SECRET not found in .env, using insecure default.")
    return signupsecret if opt == "ss" else (cookiesecret if opt == "cs" else "")
