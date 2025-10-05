import argon2
import jwt
import datetime
from bottle import Response
from typing import *
from fileinterface import myopen

ph = argon2.PasswordHasher()

def checkpass(username, password, response: Response) -> bool:
    with myopen(f"users/{username}", "r") as f:
        hashed = f.read()
    try:
        ph.verify(hashed, password)
        token = jwt.encode(
            {
                'username': username,
                'exp': datetime.datetime.now() + datetime.timedelta(days=7)  # Token expires in 7 days
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
    with myopen(f"users/{username}", "w") as f:
        f.write(hashed)
    token = jwt.encode(
        {
            'username': username,
            'exp': datetime.datetime.now() + datetime.timedelta(days=7)  # Token expires in 7 days
        },
        getsecret("cs"),
        algorithm='HS256'
    )
    response.set_cookie("token", token, httponly=True, samesite='Strict', max_age=7*24*60*60)
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
