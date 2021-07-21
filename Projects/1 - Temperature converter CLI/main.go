package main

import (
	"flag"
	"fmt"
	"strconv"
)

func CelsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32.0
}

func main() {

	tempC := flag.String("tempInCelsius", "5", "pass temperature in Celsius")

	flag.Parse()

	if *tempC == "" {
		flag.Usage()
		return
	}
	tempCelsius, err := strconv.ParseFloat(*tempC, 64)
	if err != nil {
		panic(err)
	}
	tempF := CelsiusToFahrenheit(tempCelsius)
	fmt.Println("The temperature in Fahrenheit is", tempF)

}
