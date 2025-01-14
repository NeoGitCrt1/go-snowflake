<div>
  <p align="center">
    <image src="https://www.pngkey.com/png/full/105-1052235_snowflake-png-transparent-background-snowflake-with-clear-background.png" width="250" height="250">
  </p>
  <p align="center">An Lock Free ID Generator for Golang based on Snowflake Algorithm (Twitter announced).</p>
  <p align="center">
    <a href="https://github.com/godruoyi/go-snowflake/actions?workflow=run%20tests">
      <image src="https://github.com/godruoyi/go-snowflake/workflows/run%20tests/badge.svg" alt="test">
    </a>
    <a href="https://goreportcard.com/report/github.com/godruoyi/go-snowflake">
      <image src="https://goreportcard.com/badge/github.com/godruoyi/go-snowflake" alt="Go report">
    </a>
    <a href="https://coveralls.io/repos/github/godruoyi/go-snowflake/badge.svg?branch=master">
      <image src="https://coveralls.io/repos/github/godruoyi/go-snowflake/badge.svg?branch=master" alt="Coverage status">
    </a>
  </p>
</div>

## Description

An **Lock Free** ID Generator for Golang implementation.

![file](https://images.godruoyi.com/logos/201908/13/_1565672621_LPW65Pi8cG.png)

Snowflake is a network service for generating unique ID numbers at high scale with some simple guarantees.

* The first bit is unused sign bit.
* The second part consists of a 41-bit timestamp (milliseconds) whose value is the offset of the current time relative to a certain time.
* The 10 bits machineID(5 bit workid + 5 bit datacenter id), max value is 2^10 -1 = 1023.
* The last part consists of 12 bits, its means the length of the serial number generated per millisecond per working node, a maximum of 2^12 -1 = 4095 IDs can be generated in the same millisecond.
* The binary length of 41 bits is at most 2^41 -1 millisecond = 69 years. So the snowflake algorithm can be used for up to 69 years, In order to maximize the use of the algorithm, you should specify a start time for it.

The ID generated by the snowflake algorithm is not guaranteed to be unique. For example, when two different requests enter the same machine at the same time, and the sequence generated by the node is the same, the generated ID will be duplicated.

So if you want use the snowflake algorithm to generate unique ID, You must ensure: The sequence-number generated in the same millisecond of the same node is unique.

For performance optimization
I used a concept of 'Virtual Time'

That is to say I only take real millisecond at system start. 
Every 4095 time generating snowflake, my snowflake time +1 millisecond.
This mechanism prevents time back problem.

> NOTICE: IF YOU DO CARE THE REAL TIME IN THE ID, YOU MAY NOT USE THIS

user only needs to ensure that the machine is different. You can get a unique ID.

## INSPIRED BY
### github.com/godruoyi/go-snowflake



## Feature

- ✅ Lock Free
- 🎈 Zero configuration, out of the box
- 🚀 Concurrency safety
- 🌵 Support private ip to machineid

## Installation

```shell
$ go get github.com/NeoGitCrt1/go-snowflake
```

## Usage

1. simple to use.

```go
package main

import (
    "fmt"

    "github.com/NeoGitCrt1/go-snowflake"
)

func main() {
    id := snowflake.ID()
    fmt.Println(id)
    // 1537200202186752
}
```

2. Specify the MachineID.

```go

func main() {
    snowflake.SetMachineID(1)

    // Or set private ip to machineid, testing...
    // snowflake.SetMachineID(snowflake.PrivateIPToMachineID())

    id := snowflake.ID()
    fmt.Println(id)
}
```

3. Specify start time.

```go

func main() {
    snowflake.SetStartTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC))
    id := snowflake.ID()
    fmt.Println(id)
}
```

4. Parse ID.

```go

func main() {
    id := snowflake.ID()
    sid := snowflake.ParseID(id)

    fmt.Println(sid.ID)             // 132271570944000000
    fmt.Println(sid.MachineID)      // 0
    fmt.Println(sid.Sequence)       // 0
    fmt.Println(sid.Timestamp)      // 31536000000
    fmt.Println(sid.GenerateTime()) // 2009-11-10 23:00:00 +0000 UTC
}
```

## Best practices

> ⚠️⚠️ All SetXXX method is thread-unsafe, recommended you call him in the main function.

```go

func main() {
    snowflake.SetMachineID(1) // change to your machineID
    snowflake.SetStartTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC))

    http.HandleFunc("/order", submitOrder)
    http.ListenAndServe(":8090", nil)
}

func submitOrder(w http.ResponseWriter, req *http.Request) {
    orderId := snowflake.ID()
    // save order
}
```

## License

MIT
