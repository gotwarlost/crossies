# crossies

Crossword tools and website

The tools and site that power https://crossies.us

## Design principles

* High-contrast, mobile friendly site suitable for people approaching a certain age :) 
* Web scraping logic in Go. Because, good error handling
* Vanilla JS with cherry-picked libraries for touch and drag/ drop. Keeping up with shiny frameworks for a backend programmer is hard.
* Must run on shared hosting since I already pay for it elsewhere. Go backend implemented as a FastCGI server as a result.
* Support local development with a simple HTTP server that exposes the same API without FastCGI shenanigans.
* Most operations supported by the web interface also has a CLI version for power users.



