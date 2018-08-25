package semitemp

type MachineLevelConfig struct {
	Cpu  float64
	Mem  float64
	Disk int
	P    int
	M    int
	PM   int
}

func (c *MachineLevelConfig) isEqual(v *MachineLevelConfig) bool {
	return v.Cpu == c.Cpu && v.Mem == c.Mem && v.Disk == c.Disk && v.P == c.P && v.M == c.M && v.PM == c.PM
}

func (c *MachineLevelConfig) Less(v *MachineLevelConfig) bool {
	l1 := c
	l2 := v

	if l1.Cpu < l2.Cpu {
		return true
	} else if l1.Cpu == l2.Cpu {
		if l1.Mem < l2.Mem {
			return true
		} else if l1.Mem == l2.Mem {
			if l1.Disk < l2.Disk {
				return true
			} else if l1.Disk == l2.Disk {
				if l1.P < l2.P {
					return true
				} else if l1.P == l2.P {
					if l1.M < l2.M {
						return true
					} else if l1.M == l2.M {
						if l1.PM < l2.PM {
							return true
						} else {
							return false
						}
					} else {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

type MachineLevelConfigPool struct {
	ConfigList []*MachineLevelConfig
}

func NewMachineLevelConfigPool() *MachineLevelConfigPool {
	p := &MachineLevelConfigPool{}
	return p
}

func (p *MachineLevelConfigPool) GetConfig(config *MachineLevelConfig) (result *MachineLevelConfig) {
	for _, v := range p.ConfigList {
		if v.isEqual(config) {
			return v
		}
	}

	//debugLog("MachineLevelConfigPool.GetConfig new level %v", config)

	result = &(*config)
	p.ConfigList = append(p.ConfigList, result)

	return result
}
