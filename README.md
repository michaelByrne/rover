# rover

to run the cli app, first run 

`go mod vendor` 

and then run

`go run main.go`

with possible options: `-name <rover name> -days <number of days back you want to go> -dailyMax <the max number of images to fetch per day> -camera <which camera you'd like>`

the app has a cache but it's sort of worthless because it's in-memory and each run of the app is a new instance of the app. in other words, the cache is nuked each time the app finishes

