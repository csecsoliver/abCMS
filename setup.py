import os
import platform
match platform.system():
    case "Windows":
        os.system("setup.bat")
    case "Linus":
        print("WTF are you doing?")
    case "Linux":
        print("detected linux, installing backend dependencies")
        print("You will need to have caddy and curl installed, and a wildcard domain pointing to this ip:")
        os.system("curl -4 ifconfig.me")
        os.system("sh ./setup.sh")
    case "Darwin":
        print("MacOS not yet supported: Feel free to contribute if you need.it.")
    case "Java":
        print("Why? Cannot clear the screen.")
    case _:
        print("Incompatible system, cannot clear the screen")
