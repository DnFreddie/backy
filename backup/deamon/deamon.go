package deamon

import (
	"fmt"
	"time"
)

func loopEnternely(){
	for{

	fmt.Println("Looppign")
	time.Sleep(2 * time.Second)
	}

}
