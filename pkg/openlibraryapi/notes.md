# Open Library API

## Data Model 

Books are normally "works", though sometimes they might be a specific "edition" of a work"

work - a collection of similar editions. e.g. If there is a translation or special edition of a book (work). Name of the Wind has several translation editions as well as anniversary special editions.

## api endpoints

`/works/{work id}.json` - return a work collection
`/books/{edition id}.json` - return information about a particular edition
`/isbn/{id}.json` - retuns an edition page, by searching either the isbn 10 or isbn 13 number
