package main

import (
	"database/sql"
	"log"
	"os"
	"vran/blogaggre/internal/config"
	"vran/blogaggre/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg        *config.Config
	db_queries *database.Queries
}

var mycommands = commands{container: map[string]func(*state, command) error{}}

func main() {
	mycommands.register("login", handlerLogin)
	mycommands.register("register", handlerRegister)
	mycommands.register("reset", handlerReset)
	mycommands.register("users", handlerUsers)
	mycommands.register("agg", handlerAgg)
	mycommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	mycommands.register("feeds", handlerFeeds)
	mycommands.register("follow", middlewareLoggedIn(handlerFollow))
	mycommands.register("following", middlewareLoggedIn(handlerFollowing))

	c, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// register "postgres" driver with "database/sql"
	dbURL := c.DatabaseURL
	db, err := sql.Open("postgres", dbURL)

	queries := database.New(db)

	current_state := state{
		cfg:        &c,
		db_queries: queries,
	}

	if len(os.Args) < 2 {
		log.Fatal("Missing command!")
	}

	cmd_name := os.Args[1]
	cmd_args := os.Args[2:]
	mycmd := command{
		name:      cmd_name,
		arguments: cmd_args,
	}

	if err := mycommands.run(&current_state, mycmd); err != nil {
		log.Fatal(err)
	}
}
