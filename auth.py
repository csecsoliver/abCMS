"""This module handles the authentication and data management of users

Methods:
    checkpass
    checkcookie
    createuser
    get_coins
    get_xp
    update_coins
    getsecret
    get_prefs
    save_prefs
"""

import datetime
import json
from typing import Literal

import argon2
import jwt
from bottle import Request, Response

from fileinterface import myopen
from pathlib import Path
ph = argon2.PasswordHasher()


def checkpass(username: str, password: str, response: Response) -> bool:
    """Checks the given password against the given user in the database

    Args:
        username (str): Username of the user. Will not be sanitized or checked for anything special.
        password (str): Plaintext password of the user
        response (Response): The response object to modify with the cookie

    Raises:
        Nothing, exceptions should be handled internally

    Returns:
        bool: The success of the operation
    """
    try:
        path = Path(f"users/{username}")
        if path.exists(): 
            with myopen(f"users/{username}", "r") as f:
                user_data = json.load(f)
            hashed = user_data["password"]
        else:
            raise Exception("user nonexistent")
    except:
        response.status = 401
        return False
    try:
        _ = ph.verify(hashed, password)
        token = jwt.encode(
            {
                "username": username,
                "exp": datetime.datetime.now() + datetime.timedelta(days=7),
            },
            getsecret("cs"),
            algorithm="HS256",
        )
        response.set_cookie(
            "token", token, httponly=True, samesite="Strict", max_age=7 * 24 * 60 * 60
        )
        return True
    except argon2.exceptions.VerifyMismatchError:
        response.status = 401
        return False


def checkcookie(request: Request) -> str | None:
    """Checks the jwt cookie using the cookiesecret. Returns the user

    Args:
        request (Request): The request that contains the cookie to check

    Returns:
        str | None: The username of the user or None if not valid
    """
    token = request.get_cookie("token")
    if token is None:
        return None
    try:
        payload = jwt.decode(token, getsecret("cs"), algorithms=["HS256"])
        return payload["username"]
    except (jwt.ExpiredSignatureError, jwt.InvalidTokenError):
        return None


def createuser(username: str, password: str, secret: str, response: Response) -> bool:
    """Creates a user with the given parameters

    Args:
        username (str): The username to create the user with. Will NOT be sanitized.
        password (str): The plaintext password of the user
        secret (str): The signup secret to check against the environment variables
        response (Response): The response to modify with the cookie

    Returns:
        bool: Success of the action.
    """
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
            "username": username,
            "exp": datetime.datetime.now() + datetime.timedelta(days=7),
        },
        getsecret("cs"),
        algorithm="HS256",
    )
    response.set_cookie(
        "token", token, httponly=True, samesite="Strict", max_age=7 * 24 * 60 * 60
    )
    return True


def get_coins(username: str) -> int:
    """Get the number of coins/mana of the specified user. (This is the one that can be spent)

    Args:
        username (str): Username of the user to get

    Returns:
        int: The number of currency the user has
    """
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    return user_data.get("coins", 0)


def get_xp(username: str) -> int:
    """Get the xp count of the user. (This can only change upwards programmatically)

    Args:
        username (str): Username of the user

    Returns:
        int: Experience point count
    """
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    return user_data.get("xp", 0)


def update_coins(username: str, amount: int) -> bool:
    """Changes the number of coins of the user by the specified amount

    Args:
        username (str): Username of the user
        amount (int): The amount to change by. Can be negative.

    Returns:
        bool: The success of the operation (False if the user does not have enough coins)
    """
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
    """The method to get the environment variables (not really) that store the secrets

    Args:
        opt (Literal[&quot;ss&quot;, &quot;cs&quot;]): "ss" means signup secret, "cs" mean cookie secret

    Returns:
        str: The value of the chosen secret
    """
    signupsecret = "securesignup"
    cookiesecret = "asdfkjsddhgfdzjkjsdf"
    with myopen(".env", "r") as f:
        for line in f.readlines():
            if line.startswith("SIGNUP_SECRET="):
                signupsecret = line.split("=", 1)[1].strip()
                print(f"Found SIGNUP_SECRET in .env, using that.")
            elif line.startswith("COOKIE_SECRET="):
                cookiesecret = line.split("=", 1)[1].strip()
    if signupsecret == "securesignup":
        print("SIGNUP_SECRET not found in .env, using insecure default.")
    if cookiesecret == "asdfkjsddhgfdzjkjsdf":
        print("COOKIE_SECRET not found in .env, using insecure default.")
    return signupsecret if opt == "ss" else (cookiesecret if opt == "cs" else "")


def get_prefs(username: str) -> dict[str, str]:
    """Gets the specifies user's preferences

    Args:
        username (str): Username of the user

    Returns:
        dict[str, str]: Cointains all the preferences found. Nothing is guaranteed to be in it
    """
    try:
        with myopen(f"users/{username}", "r") as f:
            user_data = json.load(f)
    except:
        return {}
    return user_data.get("prefs", {"social": ""})


def save_prefs(username: str, prefs: dict[str, str]) -> None:
    """Sets the preferences of the user

    Args:
        username (str): Username of the user
        prefs (dict[str, str]): Should contain all preferences. Any data can be stored, but research preexisting fields before adding a new one.
    """
    with myopen(f"users/{username}", "r") as f:
        user_data = json.load(f)
    user_data["prefs"] = prefs
    with myopen(f"users/{username}", "w") as f:
        json.dump(user_data, f)
