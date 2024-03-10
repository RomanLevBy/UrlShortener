package main

import (
    "fmt"
	"github.com/RomanLevBy/UrlShortener/internal/config"
)

func main() {
    //export CONFIG_PATH=./config/local.yaml
	conf := config.MustLoad()

	_ = conf

	//todo init logger

	//todo init storage

	//todo init router: chi, render

	//todo run server
}

