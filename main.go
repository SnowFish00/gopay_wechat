package main

import (
	"fmt"
	config_viper "pay/config/vipper"
)

func main() {
	fmt.Println(config_viper.Config())
}
