package discrete

import "fmt"

func ExampleStabilising() {
	const res = 1000
	driver:= NewDamped(res,0)
	driver.Step(1)
	driven:= NewStabilising(NewDamped(res,.5),0.1)
	// drive oscillators until driven is triggered
	var steps uint
	for !driven.StepChange(real(driver.State)){
		driver.Step(0)
		steps++
	}			
	fmt.Printf("#%d:%v",steps,driven)
	// Output:
	// <nil>
}