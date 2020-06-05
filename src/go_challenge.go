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

func getCPUSample() (idle, total uint64) {
    contents, err := ioutil.ReadFile("/proc/stat")
    if err != nil {
        return
    }
    lines := strings.Split(string(contents), "\n")
    for _, line := range(lines) {
        fields := strings.Fields(line)
        if fields[0] == "cpu" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                val, err := strconv.ParseUint(fields[i], 10, 64)
                if err != nil {
                    fmt.Println("Error: ", i, fields[i], err)
                }
                total += val // tally up all the numbers to get total ticks
                if i == 4 {  // idle is the 5th field in the cpu line
                    idle = val
                }
            }
            return
        }
    }
    return
}

func getCPUUsage()(float64){
    idle0, total0 := getCPUSample()
    time.Sleep(time.Second)
    idle1, total1 := getCPUSample()

    idleTicks := float64(idle1 - idle0)
    totalTicks := float64(total1 - total0)
    cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

    return cpuUsage
}

func main() {
    s1, s2, s3, s4 := getSensorData()
    cpu = getCPUUsage()
}