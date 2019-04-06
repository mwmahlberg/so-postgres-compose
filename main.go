package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	opts pg.Options
)

// User represents a user in our database
type User struct {
	Id     int64
	Name   string
	Emails []string
}

func init() {
	// We take the various configuration parameters from environment variables.
	// This makes it easy to have the docker image wrapping this application
	// to run in various environments, say with k8s config maps.
	// In a production application, you might want to use something like
	// github.com/spf13/cobra/cobra or
	// gopkg.in/alecthomas/kingpin.v2
	// instead
	opts.User = mustHaveEnv("POSTGRES_USER")
	opts.Password = mustHaveEnv("POSTGRES_PASSWORD")
	opts.Addr = fmt.Sprintf("%s:%s", mustHaveEnv("POSTGRES_HOST"), mustHaveEnv("POSTGRES_PORT"))
	opts.Database = mustHaveEnv("POSTGRES_DB")

	// We want to make sure we are connected
	opts.OnConnect = func(*pg.Conn) error {
		log.Println("Callback: Connected with database!")
		return nil
	}
}

func main() {
	log.Println("Starting up...")

	// Beware: go-pg/pg seems to lazily connect...
	db := pg.Connect(&opts)
	defer db.Close()

	// We want to make sure we actually do have a schema
	if err := createSchema(db); err != nil {
		panic(err)
	}

	// Here, you could set up your api. This is obviously just a crude example.
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(db.PoolStats())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Helper function to ensure the required parameters are existing.
func mustHaveEnv(env string) string {
	val := os.Getenv(env)
	if val == "" {
		panic(env + " must be present")
	}
	return val
}

// Helper function to create a schema
// Blatantly copied from https://github.com/go-pg/pg#look--feel
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			// Difference: We want a persistent table
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
