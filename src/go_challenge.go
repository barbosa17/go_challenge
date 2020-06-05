package main

import (
    "fmt"
    "io/ioutil"
    "math/rand"
    "strconv"
    "strings"
)

func getSensorData() (int, int, int, int) {

    rand.Seed(time.Now().UnixNano())

    sensor1 := rand.Intn(100)
    sensor2 := rand.Intn(100)
    sensor3 := rand.Intn(100)
    sensor4 := rand.Intn(100)

    return sensor1, sensor2, sensor3, sensor4
}

func main() {
    s1, s2, s3, s4 := getSensorData()
}