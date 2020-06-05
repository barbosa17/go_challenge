package main

import (
    "bufio"
    "database/sql"
    "fmt"
    "io/ioutil"
    "math/rand"
    "os"
    "runtime"
    "strconv"
    "strings"
    "time"
    _ "github.com/mattn/go-sqlite3"
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

func getMemUsage() (float64){
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    totalAlloc := m.TotalAlloc / 1024
    sysRAM := m.Sys / 1024
    var ramUsage = 100 * (float64(totalAlloc) / float64(sysRAM))
    return ramUsage
}

func storeData(timeout chan bool) {
    database, _ := sql.Open("sqlite3", "./measures.db")
    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS resources (id INTEGER PRIMARY KEY, cpu FLOAT, mem FLOAT, sensor1 INTEGER, sensor2 INTEGER, sensor3 INTEGER, sensor4 INTEGER)")
    statement.Exec()
    statement, _ = database.Prepare("INSERT INTO resources (cpu, mem, sensor1, sensor2, sensor3, sensor4) VALUES (?, ?, ?, ?, ?, ?)")
    s1,s2,s3,s4 := getSensorData()
    statement.Exec(getCPUUsage(), getMemUsage(), s1,s2,s3,s4)
    timeout <- true
}

func readData(n int, props string) {
    trimed_props := strings.ReplaceAll(props, ",", "")
    fixed_props := strings.ReplaceAll(trimed_props, " ", ", ")

    database, err := sql.Open("sqlite3", "./measures.db")

    rows, err := database.Query("SELECT " + fixed_props + " FROM resources ORDER BY id DESC LIMIT " + strconv.Itoa(n))
    if err != nil {
        fmt.Printf("Error: %v", err)
        return
    }
    columns, err := rows.Columns()
    if err != nil {
        fmt.Printf("Error: %v", err)
        return
    }
    colNum := len(columns)

    var values = make([]interface{}, colNum)
    for i, _ := range values {
        var ii interface{}
        values[i] = &ii
    }

    for rows.Next() {
        err := rows.Scan(values...)
        if err != nil {
            fmt.Printf("Error: %v", err)
            return
        }
        for i, colName := range columns {
            var raw_value = *(values[i].(*interface{}))
            fmt.Println(colName,raw_value)

        }
    }
}

// TODO: Revisit this function
// Implement this behaviour in the readData()
func readDataAVG(n int, props string) {
    trimed_props := strings.ReplaceAll(props, ",", "")

    database, _ := sql.Open("sqlite3", "./measures.db")

    avg_props:= ""
    words := strings.Fields(trimed_props)
    for i, elem := range words {
        if (i == 0){
            avg_props = "AVG(" + string(elem) + ") " + avg_props
        }else {
            avg_props = "AVG(" + string(elem) + "), " + avg_props
        }
    }
    rows, err := database.Query("SELECT " + avg_props + " FROM resources")

    if err != nil {
        fmt.Printf("Error: %v", err)
        return
    }
    columns, err := rows.Columns()
    if err != nil {
        fmt.Printf("Error: %v", err)
        return
    }
    colNum := len(columns)

    var values = make([]interface{}, colNum)
    for i, _ := range values {
        var ii interface{}
        values[i] = &ii
    }

    for rows.Next() {
        err := rows.Scan(values...)
        if err != nil {
            fmt.Printf("Error: %v", err)
            return
        }
        for i, colName := range columns {
            var raw_value = *(values[i].(*interface{}))
            fmt.Println(colName,raw_value)

        }
    }
}

func menu(main chan bool) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("\n\nUbiWhere challenge Golang\n\n")
    fmt.Print("Please choose one of the following options:\n")
    fmt.Print("1: for get an amount all metrics\n")
    fmt.Print("2: for specify what metrics you want\n")
    fmt.Print("3: for an average of one of more metrics\n")
    fmt.Print("Other integer value will quit the application.\n")
    menu, _ := reader.ReadString('\n')
    menu = strings.Replace(menu, "\n", "", -1)
    var option, err = strconv.Atoi(menu)
    if (err == nil) {
        switch (option) {
            case 1:
                fmt.Print("Enter the amont of samples you want to get: ")
                n_metrics, _ := reader.ReadString('\n')
                n_metrics = strings.Replace(n_metrics, "\n", "", -1)
                var samples, err = strconv.Atoi(n_metrics)
                if (err == nil) {
                    readData(samples, "*")
                } else {
                    fmt.Print("Please use just integer values \n")
                }
            case 2:
                fmt.Print("Enter the amont of samples you want to get: ")
                n_metrics, _ := reader.ReadString('\n')
                n_metrics = strings.Replace(n_metrics, "\n", "", -1)
                var samples, err = strconv.Atoi(n_metrics)
                if (err == nil) {
                    fmt.Print("Please specify the property name (sensor[1.4], cpu or ram): ")
                    properties, _ := reader.ReadString('\n')
                    readData(samples, properties)
                } else {
                    fmt.Print("Please use just integer values \n")
                }
            case 3:
                fmt.Print("Will do the average \n")
                fmt.Print("Please specify the property name (sensor[1.4], cpu or ram): ")
                properties, _ := reader.ReadString('\n')
                readDataAVG(0, properties)
            default:
                return
        }
    } else {
        fmt.Print("Please use just integer values \n")
    }
    main <- true
}

func main() {
    timeout := make(chan bool)
    main := make(chan bool)
    for {
        select {
        case <-time.After(time.Second):
            go storeData(timeout)
            <- timeout
        }
        go menu(main)
        <-main
        //end loop
    }
    //end main
}