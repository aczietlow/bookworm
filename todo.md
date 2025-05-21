## Underlying API
- [x] Add command to search for books by title
- [x] Add cache
- [x] Completely refactor out legacy code 
- [x] Add openlibrary data objects for works and editions
- [ ] openlibrary URL arguments need to all be upper case. 
- [ ] Update http calls for OL objects to check the cache first
  - decide if you want to cache api calls, or aggregated data objects
- [ ] Add command to fetch raw json from api to demo tree viewer
- [ ] Add command to add books to reading list


## TUI
- [x] bootstrap 2 table layout
- [x] Display the out of commands in right talbe pane 
- [ ] Refactor PoC tui logic from main() into appropraite locations
- [ ] Add f(x) to pull data from result pane to another f(x)
  - e.g. Using the search function, search for a book by title, select the desired book, then "insepct" that book to fetch additional information about said book

## problems for later 
- [ ] build a factory that allows for easier registering of command 
- [ ] add architecture.md doc
- [ ] Add a spinner when fetching data from the api 
