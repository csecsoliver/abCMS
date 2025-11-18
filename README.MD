File attributions:
filestore/copyparty.py is made by https://github.com/9001 and is downloadable from https://copyparty.eu/py

# abCMS
This has been my ongoing project throughout siege. This week I made good progress rewriting the entire thing in moonscript (lua) with the lapis framework. Image posts work now, post colors do not

I have done this mostly to create a better framework for doing things properly.
I use sqlite for the db, copyparty for free space display and possibly image serving in the future.

If you were ever brave enough to host an instance of this, you can now migrate over!
1. You're gonna need two shell sessions, one for copyparty and one for lapis.
2. In shell #1 clone the repo into an empty dir, and run setup.sh
3. If you have the data from a previous version, copy the folders into the current repo dir and run migrate.py.
4. in shell #1 run 'lapis serve'
5. In shell #2 go to the repo dir and run setup_files.sh
6. Congrats, you now have a working site on port 8080, which you can expose using a reverse proxy. Do not expose port 3923 in the current version.


## Repo structure
- The releases only have the python versions as I have yet to figure out packaging. 
- The v3-legacy-python branch contains the code for the python version
- The v4-rewrite branch is not relevant


