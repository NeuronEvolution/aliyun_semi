package schedule

import "fmt"

type MachineCollection struct {
	Map       map[int]*Machine
	List      []*Machine
	ListCount int
}

func NewMachineCollection() *MachineCollection {
	c := &MachineCollection{}
	c.Map = make(map[int]*Machine)
	c.List = make([]*Machine, MaxMachineId+1)

	return c
}

func (c *MachineCollection) debugValidation() {
	for i := 0; i < c.ListCount; i++ {
		if c.List[i] == nil {
			panic(fmt.Errorf("MachineCollection.debugValidation c.List[%d]==nil,%d", i, c.ListCount))
		}
	}
}

func (c *MachineCollection) Add(m *Machine) {
	//debugLog("MachineCollection.Add %d", m.MachineId)

	if DebugEnabled {
		_, has := c.Map[m.MachineId]
		if has {
			panic(fmt.Errorf("MachineCollection.Add %d exists", m.MachineId))
		}
	}

	c.Map[m.MachineId] = m
	c.List[c.ListCount] = m
	c.ListCount++

	if DebugEnabled {
		c.debugValidation()
	}
}

func (c *MachineCollection) Remove(machineId int) {
	//debugLog("MachineCollection.Remove %d", machineId)

	if DebugEnabled {
		_, has := c.Map[machineId]
		if !has {
			panic(fmt.Errorf("MachineCollection.Add %d not exists", machineId))
		}
	}

	delete(c.Map, machineId)

	for i := 0; i < c.ListCount; i++ {
		v := c.List[i]
		if v.MachineId == machineId {
			c.List[i] = nil
			if c.ListCount > 1 && i < c.ListCount-1 {
				c.List[i] = c.List[c.ListCount-1]
			}
			c.ListCount--
		}
	}

	if DebugEnabled {
		c.debugValidation()
	}
}

func (c *MachineCollection) Peek() (m *Machine) {
	//debugLog("MachineCollection.Peek ListCount=%d", c.ListCount)
	if c.ListCount == 0 {
		return nil
	}

	return c.List[c.ListCount-1]
}

func (c *MachineCollection) Has(machinedId int) bool {
	_, has := c.Map[machinedId]
	return has
}

func (c *MachineCollection) First() *Machine {
	if c.ListCount == 0 {
		return nil
	}

	return c.List[0]
}

func (c *MachineCollection) Last() *Machine {
	if c.ListCount == 0 {
		return nil
	}

	return c.List[c.ListCount-1]
}
