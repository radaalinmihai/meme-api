### MIGRATIONS
To migrate the schemas into your database use the next command:

`migrate -path db/migrations -database "mysql://<username>:<password>@tcp(<your ip | localhost>:3306)/meme" -verbose up`

To revert changes and delete tables do:

`migrate -path db/migrations -database "mysql://<username>:<password>@tcp(<your ip | localhost>:3306)/meme" -verbose down`

### NOTE
You need to have installed scoop in order to use migrate

You need to already have a 'meme' database created in order to migrate the schemas

Migration library:

https://github.com/golang-migrate/migrate