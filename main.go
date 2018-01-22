package main

import (
	"fmt"
	"os"

	"github.com/ahmedalhulaibi/go-graphqlator-cli/cmd"

	_ "github.com/ahmedalhulaibi/substance/providers/mysqlsubstance"
	_ "github.com/ahmedalhulaibi/substance/providers/pgsqlsubstance"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
