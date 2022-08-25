package bootstrap

import (
	"skframe/pkg/config"
	"skframe/pkg/elastic"
)

func ConnectElastic() {
	elastic.ConnectElastic(config.GetString("es.url"))
}
