import argon2
import jwt
import datetime
from bottle import Response
from typing import *
from fileinterface import myopen
import json

ph = argon2.PasswordHasher()

def checkpass(username, password, response: Response) -> bool:
    try:
        with myopen(f"users/{username}", "r") as f:
            user_data = json.load(f)
        hashed = user_data["password"]
    except (FileNotFoundError, KeyError):
        response.status = 401
        return False
    try:
        ph.verify(hashed, password)
        token = jwt.encode(
            {
                'username': username,
                'exp': datetime.datetime.now() + datetime.timedelta(days=7)
            },
            getsecret("cs"),
            algorithm='HS256'
        )
        response.set_cookie("token", token, httponly=True, samesite='Strict', max_age=7*24*60*60)
        return True
    except argon2.exceptions.VerifyMismatchError:
        response.status = 401
        return False
    
def checkcookie(request) -> Union[str, None]:
    token = request.get_cookie("token")
    if token is None:
        return None
    try:
        payload = jwt.decode(token, getsecret("cs"), algorithms=['HS256'])
        return payload['username']
    except (jwt.ExpiredSignatureError, jwt.InvalidTokenError):
        return None

def createuser(username, password, secret, response: Response) -> bool:
    if secret != getsecret("ss"):
        return False
    try:
        with open(f"users/{username}", "r") as f:
            pass
        return False
    except FileNotFoundError:
        pass
    hashed = ph.hash(password)
    user_data = {"password": hashed, "coins": 0}
    with myopen(f"users/{username}", "w") as f:
        json.dump(user_data, f)
    token = jwt.encode(
        {
            'username': username,
            'exp': datetime.datetime.now() + datetime.timedelta(days=7)
        },
        getsecret("cs"),
        algorithm='HS256'
    )
    response.set_cookie("token", token, httponly=True, samesite='Strict', max_age=7*24*60*60)
    return True

def get_coins(username: str) -> int:
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    return user_data.get("coins", 0)

def get_xp(username: str) -> int:
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    return user_data.get("xp", 0)

def update_coins(username: str, amount: int) -> bool:
    
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    user_data["coins"] = user_data.get("coins", 0) + amount
    if user_data["coins"] < 0:
        return False
    if amount > 0:
        user_data["xp"] = user_data.get("xp", 0) + amount
    with myopen(f"users/{username}", "w") as f:
        json.dump(user_data, f)
    return True

def getsecret(opt: Literal["ss", "cs"]) -> str:
    signupsecret = "securesignup"
    cookiesecret = "asdfkjsddhgfdzjkjsdf"
    with myopen(".env", "r") as f:
        for line in f.readlines():
            if line.startswith("SIGNUP_SECRET="):
                signupsecret = line[len("SIGNUP_SECRET="):].strip()
            elif line.startswith("COOKIE_SECRET="):
                cookiesecret = line[len("COOKIE_SECRET="):].strip()
    if signupsecret == "securesignup":
        print("SIGNUP_SECRET not found in .env, using insecure default.")
    if cookiesecret == "asdfkjsddhgfdzjkjsdf":
        print("COOKIE_SECRET not found in .env, using insecure default.")
    return signupsecret if opt == "ss" else (cookiesecret if opt == "cs" else "")

def get_prefs(username: str) -> dict:
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    return user_data.get("prefs", {"social": ""})

def save_prefs(username: str, prefs: dict) -> None:
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    user_data["prefs"] = prefs
    with myopen(f"users/{username}", "w") as f:
        json.dump(user_data, f)