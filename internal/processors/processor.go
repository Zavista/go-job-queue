package processors

/*
Processor is the interface that all processes must implement to handle the work logic (fix comments later)
*/
type Processor interface {
	Process() (string, error)
	Type() string
}
