# jeopardy
Simple Jeopardy server that exposes an API for applications to obtain unique game data on each call.
A custom mySQL database is queried each time data is requested, currently accessed as a local file.

**Instructions**

Download the pre-compiled and zipped exectuable with data from the releases section.
If you want to compile it yourself, you'll need the GNU Compiler (you can get it [here](https://sourceforge.net/projects/mingw-w64/files/mingw-w64/mingw-w64-release/) for Windows)

Once you have it running, navigate to http://localhost:8080/ for information or http://localhost:8080/jeopardy for unique game data.


*Credits*

Original DB from github.com/jwolle1/jeopardy_clue_dataset

All data is property of Jeopardy Productions, Inc. and protected under law. I am not affiliated with the show. Please don't use the data to make a public-facing web site, app, or any other product.

Thanks to Code Louisville for their resources.