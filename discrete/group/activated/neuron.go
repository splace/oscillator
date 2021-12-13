package activated

import oscillator "../../../../oscillator/stepped"

//import "log"

// Neuron's are Triggering oscillators with a Threshold that changes.
// uses Different Threshold Adjusters depending on if currently Triggered.
type Neuron struct {
	*Triggering 
	sensitiser, desensitiser Adjuster
}

type Adjuster interface{
	Adjust(*Threshold)
}

type Geometric float64

func (a Geometric) Adjust(t *Threshold){
	(*t)*=Threshold(a)
}

type LimitedGeometric struct{
	Geometric
	Threshold
}

func (a LimitedGeometric) Adjust(t *Threshold){
	if *t>a.Threshold{
		(*t)=a.Threshold
		return
	}
	(*t)*=Threshold(a.Geometric)
}

// make a Neuron with fixed Geometric Threshold Adjusters.
// optionally different (de)sensitiser.
func NewNeuron(sc uint,damping float64,t Threshold,f float64 ,fs ...float64) *Neuron{
	nt:=NewTriggering(oscillator.NewDrivenAndDamped(sc,damping), t)
	if len(fs)==0{
		if f>1.0 {return nil}
		return &Neuron{nt,Geometric(1-f/float64(sc)),Geometric(1+f/float64(sc))} 
	}
	if f>1.0 || fs[0]>1 {return nil}
	return &Neuron{nt,Geometric(1-f/float64(sc)),Geometric(1+fs[0]/float64(sc))}
}

func (n *Neuron) StepChange(d float64) bool {
//	log.Print(h,h.Triggering)
	n.Step(d)
	if n.Triggered {
		n.desensitiser.Adjust(&n.Threshold)
//		n.Threshold*=float64(n.desensitiser)  //TODO squared!!
		if !n.Threshold.Exceeded(n.Abs()){
//		if n.DrivenAndDamped.Abs2() < n.Threshold2 {
			n.Triggered = Triggered(false)
			return true
		}
	} else {
		n.sensitiser.Adjust(&n.Threshold)
//		n.Threshold=1-(1-n.Threshold)*float64(n.sensitiser)

//		if n.DrivenAndDamped.Abs2() > n.Threshold2 {
		if n.Threshold.Exceeded(n.Abs()){
			n.Triggered = Triggered(true)
			return true
		}
	}
	return false
}

