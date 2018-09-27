package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/garsue/sparql/sql"
	"github.com/jinzhu/gorm"
	"gopkg.in/urfave/cli.v1"
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

	db, err := gorm.Open("sparql", "https://api.hojin-info.go.jp/sparql")
	if err != nil {
		return err
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	var hojins []struct {
		ID   uint64
		Name string
	}

	//noinspection SqlNoDataSourceInspection
	if err = db.Raw(`
PREFIX hj: <http://hojin-info.go.jp/ns/domain/biz/1#>
PREFIX ic: <http://imi.go.jp/ns/core/rdf#>
	
SELECT ?id,?name FROM <http://hojin-info.go.jp/graph/hojin>
WHERE {
	?s hj:法人基本情報 ?n .
	?n ic:ID/ic:識別値 ?id .
	?n ic:名称/ic:表記 ?name .
	FILTER regex(?name, $1)
} LIMIT 100`, name).Find(&hojins).Error; err != nil {
		return err
	}

	for _, h := range hojins {
		fmt.Println(h.ID, h.Name)
	}

	return nil
}
