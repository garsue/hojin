package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/garsue/sparql"
	"github.com/jinzhu/gorm"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "hojin"
	app.Usage = "command line tool for hojin-info"
	app.Action = func(c *cli.Context) error {
		return search(c.Args().First())
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func search(name string) (err error) {
	if name == "" {
		return errors.New("specify office name")
	}

	dsn := "http://api.hojin-info.go.jp/sparql"
	sql.Register("sparql", sparql.NewConnector(dsn).Driver())

	db, err := gorm.Open("sparql", dsn)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	var hojins []struct {
		ID          uint64
		Name        string
		EmergedAt   string
		Description string
		AddressType string
	}

	//noinspection SqlNoDataSourceInspection
	if err = db.Raw(`
PREFIX hj: <http://hojin-info.go.jp/ns/domain/biz/1#>
PREFIX ic: <http://imi.go.jp/ns/core/rdf#>
	
SELECT * FROM <http://hojin-info.go.jp/graph/hojin>
WHERE {
	?s hj:法人基本情報 ?n .
	?n ic:ID/ic:識別値 ?id ;
	ic:名称/ic:表記 ?name ;
	ic:活動状況/ic:発生日/ic:標準型日時 ?emerged_at ;
	ic:活動状況/ic:説明 ?description ;
	ic:住所/ic:種別 ?address_type ;
	ic:住所/ic:郵便番号 ?zip_code.
	FILTER regex(?name, $1)
} LIMIT 100`, name).Find(&hojins).Error; err != nil {
		return err
	}

	for _, h := range hojins {
		fmt.Printf("%+v\n", h)
	}

	return nil
}
