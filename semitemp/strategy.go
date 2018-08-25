package semitemp

import "fmt"

type Strategy interface {
	Name() string
	AddInstanceList(instanceList []*Instance) (err error)
}

type defaultStrategy struct {
}

func (s *defaultStrategy) Name() string {
	return "defaultStrategy"
}

func (s *defaultStrategy) AddInstanceList(instanceList []*Instance) (err error) {
	fmt.Println("defaultStrategy.AddInstanceList")
	return nil
}
