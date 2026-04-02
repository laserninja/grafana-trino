package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/sqlds/v4"
	"github.com/grafana/trino/pkg/trino"
)

func main() {
	ds := sqlds.NewDatasource(&trino.Datasource{})
	if err := datasource.Manage("trino-datasource", ds.NewDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
