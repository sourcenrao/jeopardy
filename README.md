# jeopardy
Simple Jeopardy server that exposes an API for applications to obtain unique game data on each call.

A custom mySQL database is queried each time data is requested, currently accessed as a local file.

Start the server with go run and navigate to http://localhost:8080/ for information or http://localhost:8080/jeopardy for game data.

Original DB from github.com/jwolle1/jeopardy_clue_dataset

All data is property of Jeopardy Productions, Inc. and protected under law. I am not affiliated with the show. Please don't use the data to make a public-facing web site, app, or any other product.