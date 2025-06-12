## Underlying API
- [x] Add command to search for books by title
- [x] Add cache
- [x] Completely refactor out legacy code 
- [x] Add openlibrary data objects for works and editions
- [ ] openlibrary URL arguments need to all be upper case. 
- [ ] Update http calls for OL objects to check the cache first
  - decide if you want to cache api calls, or aggregated data objects
- [ ] Add command to fetch raw json from api to demo tree viewer <--
- [ ] Add command to add books to reading list
- [x] api errors
  - [x] validate the correct ID is sent
    - was assuming json fields would be populated.
    - don't look for a slice index before checking the len of the slice
  - [x] the search term "the fault" causes a panic
    - the api results can return a works object with no author_name

## TUI
- [x] bootstrap 2 table layout
- [x] Display the out of commands in right talbe pane 
- [x] Refactor PoC tui logic from main() into appropraite locations
- [ ] unbreak input capture that allows user to esc from result pane back to view pane
- [ ] refactor command registry
  - need structs defined for each commandView since each will have its own udpate() and updateResult()
  - confirm results will actually be different
  - or if a generic we could just extend the tview.Primitives themselves to include update() functions
- [ ] Add f(x) to pull data from result pane to another f(x)
  - e.g. Using the search function, search for a book by title, select the desired book, then "insepct" that book to fetch additional information about said book
- [ ] Update how the result pane gets drawn
  - Need to allow commands to return more †han just strings e.g. hierarchical json data to get displayed as treeview

## problems for later 
- [ ] build a factory that allows for easier registering of command 
- [x] add architecture.md doc
- [ ] Add a spinner when fetching data from the api 
