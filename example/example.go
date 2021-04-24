package main

import (
	"fmt"

	"github.com/NeoGitCrt1/go-snowflake"
)

func main() {
	id := snowflake.ID()
	fmt.Println(id)

	sid := snowflake.ParseID(id)
	// SID {
	//     Sequence: 0
	//     MachineID: 0
	//     Timestamp: x
	//     ID: x
	// }
	fmt.Println(sid)
}
