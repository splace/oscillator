package discrete

import "fmt"

func ExampleTobytesConvert(){
	fmt.Println(toBytes([]uint16{1,2}))
	fmt.Println(toBytes([]uint16{256,257}))
	fmt.Println(toBytes([]uint32{256,257}))
	fmt.Println(toBytes([]uint64{256,257}))
	// Output:
	// [1 0 2 0]
	// [0 1 1 1]
	// [0 1 0 0 1 1 0 0]
	// [0 1 0 0 0 0 0 0 1 1 0 0 0 0 0 0]
}

