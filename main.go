package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"

	_ "github.com/garsue/sparql/sql"
)

func main() {
	app := cli.NewApp()

	app.Name = "hojin"
	app.Usage = "command line tool for hojin-info"
	app.Action = func(c *cli.Context) error {
		return search(context.Background(), c.Args().First())
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func search(ctx context.Context, name string) (err error) {
	if name == "" {
		return errors.New("specify office name")
	}

	db, err := sql.Open("sparql", "https://api.hojin-info.go.jp/sparql")
	if err != nil {
		return err
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	//noinspection SqlNoDataSourceInspection
	rows, err := db.Query(`
PREFIX hj: <http://hojin-info.go.jp/ns/domain/biz/1#>
PREFIX ic: <http://imi.go.jp/ns/core/rdf#>
	
SELECT ?v,?i FROM <http://hojin-info.go.jp/graph/hojin>
WHERE {
	?s hj:法人基本情報 ?n .
	?n ic:名称/ic:表記 ?v ;
	   ic:ID/ic:識別値 ?i .
	FILTER regex(?v, $1)
} LIMIT 100`, name)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := rows.Close(); err2 != nil {
			err = err2
		}
	}()

	for rows.Next() {
		var v string
		var i string
		if err := rows.Scan(&v, &i); err != nil {
			return err
		}
		fmt.Println(v, i)
	}

	return nil
}
