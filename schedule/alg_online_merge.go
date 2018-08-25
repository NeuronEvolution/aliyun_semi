package schedule

import (
	"fmt"
)

type OnlineMerge struct {
	R *ResourceManagement

	FinalMachineMap [][]int
	FinalDeployMap  []int

	MachineList         []*Machine
	MachineMap          []*Machine
	InstanceList        []*Instance
	InstanceMap         []*Instance
	DeployMap           []*Machine
	DeployedMachineList []*Machine
	DeployedMachineMap  []*Machine
	FreeMachineList     []*Machine
	FreeMachineMap      []*Machine

	MoveCommands []*InstanceMoveCommand
}

func NewOnlineMerge(r *ResourceManagement) *OnlineMerge {
	m := &OnlineMerge{}
	m.R = r
	return m
}

func (o *OnlineMerge) init() {
	//复制最终状态，不影响外部的调度
	o.copyFinalStatus()

	//创建机器及实例
	o.MachineList, o.MachineMap = o.R.createMachines()
	o.InstanceList, o.InstanceMap = o.R.createInstances()

	//部署初始状态
	o.DeployMap = make([]*Machine, o.R.MaxInstanceId+1)
	for _, config := range o.R.InstanceDeployConfigList {
		if config == nil {
			continue
		}

		instance := o.InstanceMap[config.InstanceId]
		m := o.MachineMap[config.MachineId]
		m.AddInstance(instance)
		o.DeployMap[instance.InstanceId] = m
	}

	//交换机器，固定住大的实例
	o.fixMachines()

	//区分部署和空闲机器
	o.DeployedMachineMap = make([]*Machine, o.R.MaxMachineId+1)
	o.FreeMachineMap = make([]*Machine, o.R.MaxMachineId+1)
	for machineId, m := range o.FinalMachineMap {
		machine := o.MachineMap[machineId]
		if machine == nil {
			continue
		}

		if len(m) == 0 {
			o.FreeMachineMap[machineId] = machine
			o.FreeMachineList = append(o.FreeMachineList, machine)
		} else {
			o.DeployedMachineMap[machineId] = machine
			o.DeployedMachineList = append(o.DeployedMachineList, machine)
		}
	}
}

func (o *OnlineMerge) copyFinalStatus() {
	//复制机器部署的实例列表
	o.FinalMachineMap = make([][]int, len(o.R.MachineMap))
	for i, m := range o.R.MachineMap {
		if m != nil && m.InstanceListCount > 0 {
			for _, instance := range m.InstanceList[:m.InstanceListCount] {
				o.FinalMachineMap[i] = append(o.FinalMachineMap[i], instance.InstanceId)
			}
		}
	}

	//复制实例部署机器的映射
	o.FinalDeployMap = make([]int, len(o.R.DeployMap))
	for i, v := range o.R.DeployMap {
		if v != nil {
			o.FinalDeployMap[i] = v.MachineId
		}
	}
}

func (o *OnlineMerge) swapMachine(m1 int, m2 int) {
	temp := o.FinalMachineMap[m2]
	o.FinalMachineMap[m2] = o.FinalMachineMap[m1]
	o.FinalMachineMap[m1] = temp

	if o.FinalMachineMap[m1] != nil {
		for _, instanceId := range o.FinalMachineMap[m1] {
			o.FinalDeployMap[instanceId] = m1
		}
	}

	if o.FinalMachineMap[m2] != nil {
		for _, instanceId := range o.FinalMachineMap[m2] {
			o.FinalDeployMap[instanceId] = m2
		}
	}
}

func (o *OnlineMerge) fixMachines() {
	o.R.log("OnlineMerge.fixMachines start\n")

	SortInstanceByTotalMaxLowWithInference(o.InstanceList, 32)

	fixedMachines := make(map[int]int)
	for _, instance := range o.InstanceList {
		config := o.R.InstanceDeployConfigMap[instance.InstanceId]
		currentMachineId := o.FinalDeployMap[instance.InstanceId]

		//实例当前部署的机器已被固定，不做映射
		_, fixed := fixedMachines[currentMachineId]
		if fixed {
			continue
		}

		//实例刚好部署在目标机器上，直接固定该机器
		if currentMachineId == config.MachineId {
			fixedMachines[currentMachineId] = currentMachineId
			continue
		}

		//若目标机器未被固定，将实例交换到目标机器并固定
		_, fixed = fixedMachines[config.MachineId]
		if fixed {
			continue
		}
		//fmt.Println("fixMachines", config.MachineId, currentMachineId, instance.InstanceId)
		if o.R.MachineConfigMap[config.MachineId].Cpu != o.R.MachineConfigMap[currentMachineId].Cpu {
			continue
		}
		o.swapMachine(currentMachineId, config.MachineId)
		fixedMachines[config.MachineId] = config.MachineId
	}

	o.R.log("OnlineMerge.fixMachines fixedMachines=%d\n", len(fixedMachines))
}

func (o *OnlineMerge) roundFirst() {
	o.R.log("OnlineMerge.roundFirst start\n")

	SortInstanceByTotalMaxLowWithInference(o.InstanceList, 32)

	//迁移后留下的幽灵实例
	ghosts := make([]*Instance, 0)
	ghostsDeploy := make([]*Machine, 0)

	//尝试迁移其他实例到目标机器，若失败迁移到剩余机器
	//若再失败，尝试将大实例的目标机器迁移开，下一轮迁入
	moveAlready := 0
	moveSuccess := 0
	moveTemp := 0
	moveRest := 0
	moveFreezing := 0
	for i, instance := range o.InstanceList {
		if i > 0 && i%10000 == 0 {
			o.R.log("OnlineMerge.roundFirst 1 %d\n", i)
		}

		currentMachine := o.DeployMap[instance.InstanceId]
		targetMachineId := o.FinalDeployMap[instance.InstanceId]

		//已经部署，直接跳过
		if currentMachine.MachineId == targetMachineId {
			moveAlready++
			continue
		}

		//当前已在临时机器中，保持不动
		if o.FreeMachineMap[currentMachine.MachineId] != nil {
			moveFreezing++
			continue
		}

		targetMachine := o.MachineMap[targetMachineId]
		if targetMachine.ConstraintCheck(instance, 1) {
			//迁移到目标机器
			currentMachine.RemoveInstance(instance.InstanceId)
			targetMachine.AddInstance(instance)
			o.DeployMap[instance.InstanceId] = targetMachine

			//留下幽灵
			ghost := instance.CreateGhost()
			currentMachine.AddInstance(ghost)
			ghosts = append(ghosts, ghost)
			ghostsDeploy = append(ghostsDeploy, currentMachine)

			//纪录指令
			o.MoveCommands = append(o.MoveCommands, &InstanceMoveCommand{
				Round:      1,
				InstanceId: instance.InstanceId,
				MachineId:  targetMachine.MachineId,
			})

			moveSuccess++
		} else {
			//迁移到剩余机器
			moved := false
			for _, freeMachine := range o.FreeMachineList {
				if freeMachine.ConstraintCheck(instance, 1) {
					currentMachine.RemoveInstance(instance.InstanceId)
					freeMachine.AddInstance(instance)
					o.DeployMap[instance.InstanceId] = freeMachine

					//留下幽灵
					ghost := instance.CreateGhost()
					currentMachine.AddInstance(ghost)
					ghosts = append(ghosts, ghost)
					ghostsDeploy = append(ghostsDeploy, currentMachine)

					//纪录指令
					o.MoveCommands = append(o.MoveCommands, &InstanceMoveCommand{
						Round:      1,
						InstanceId: instance.InstanceId,
						MachineId:  freeMachine.MachineId,
					})

					moved = true
					moveTemp++
					break
				}
			}
			if !moved { //todo
				moveRest++
			}
		}
	}

	o.R.log("OnlineMerge.roundFirst 1,machines=%d,deployed=%d,free=%d,"+
		"moveAlready=%d,moveSuccess=%d,moveTemp=%d，moveFreezing＝%d,moveRest=%d\n",
		len(o.MachineList), len(o.DeployedMachineList), len(o.FreeMachineList),
		moveAlready, moveSuccess, moveTemp, moveFreezing, moveRest)

	//根据第二轮状态尝试将影响下轮回放的实例移开
	nextMachineMap := make([]*Machine, o.R.MaxMachineId+1)
	for i, m := range o.MachineMap {
		if m != nil {
			nextMachineMap[i] = NewMachine(m.R, m.MachineId, m.Config)
		}
	}
	for instanceId, machineId := range o.FinalDeployMap {
		if machineId > 0 {
			//排除掉在临时机器中的实例
			currentMachine := o.DeployMap[instanceId]
			if o.FreeMachineMap[currentMachine.MachineId] != nil {
				//continue#注释掉这里限制下轮的移动，避免第一轮过度移动
			}

			nextMachineMap[machineId].AddInstance(o.InstanceMap[instanceId])
		}
	}

	moveAlready = 0
	moveFreezing = 0
	moveKeep := 0
	moveOther := 0
	moveRest = 0
	lastFitPos := 0
	for i, instance := range o.InstanceList {
		if i > 0 && i%10000 == 0 {
			o.R.log("OnlineMerge.roundFirst 2 %d\n", i)
		}

		currentMachine := o.DeployMap[instance.InstanceId]
		targetMachineId := o.FinalDeployMap[instance.InstanceId]

		//已经部署，直接跳过
		if currentMachine.MachineId == targetMachineId {
			moveAlready++
			continue
		}

		//在临时机器中，本轮不移动
		if o.FreeMachineMap[currentMachine.MachineId] != nil {
			moveFreezing++
			continue
		}

		//若不可保持位置，迁移走
		nextCurrentMachine := nextMachineMap[currentMachine.MachineId]
		if nextCurrentMachine.ConstraintCheck(instance, 1) {
			//可以保持位置
			nextCurrentMachine.AddInstance(instance)
			moveKeep++
		} else {
			//迁移走
			moved := false
			for fitOffset := 1; fitOffset <= len(o.DeployedMachineList); fitOffset++ {
				pos := lastFitPos + fitOffset
				if pos == len(o.DeployedMachineList) {
					pos = 0
				}

				deployMachine := o.DeployedMachineList[pos]
				if deployMachine == currentMachine {
					continue
				}

				if !deployMachine.ConstraintCheck(instance, 1) {
					continue
				}

				nextDeployMachine := nextMachineMap[deployMachine.MachineId]
				if !nextDeployMachine.ConstraintCheck(instance, 1) {
					continue
				}

				//迁移实例
				currentMachine.RemoveInstance(instance.InstanceId)
				deployMachine.AddInstance(instance)
				o.DeployMap[instance.InstanceId] = deployMachine

				//更新最终状态
				nextDeployMachine.AddInstance(instance)

				//留下幽灵
				ghost := instance.CreateGhost()
				currentMachine.AddInstance(ghost)
				ghosts = append(ghosts, ghost)
				ghostsDeploy = append(ghostsDeploy, currentMachine)

				//纪录指令
				o.MoveCommands = append(o.MoveCommands, &InstanceMoveCommand{
					Round:      1,
					InstanceId: instance.InstanceId,
					MachineId:  deployMachine.MachineId,
				})

				moveOther++
				moved = true
				break
			}
			if !moved {
				moveRest++
			}
		}
	}

	//删除幽灵实例
	for i, ghost := range ghosts {
		ghostsDeploy[i].RemoveInstance(ghost.InstanceId)
	}

	o.R.log("OnlineMerge.roundFirst 2,machines=%d,deployed=%d,free=%d,"+
		"moveAlready=%d,moveSuccess=%d,moveTemp=%d，moveFreezing＝%d,moveKeep=%d,moveOther=%d,moveRest=%d\n",
		len(o.MachineList), len(o.DeployedMachineList), len(o.FreeMachineList),
		moveAlready, moveSuccess, moveTemp, moveFreezing, moveKeep, moveOther, moveRest)
}

func (o *OnlineMerge) roundSecond() {
	o.R.log("OnlineMerge.roundSecond start\n")

	SortInstanceByTotalMaxLowWithInference(o.InstanceList, 32)

	//迁移后留下的幽灵实例
	ghosts := make([]*Instance, 0)
	ghostsDeploy := make([]*Machine, 0)

	//首先将可以移到目标机器的移过去
	moveAlready := 0
	moveFreezing := 0
	moveSuccess := 0
	moveRest := 0
	for i, instance := range o.InstanceList {
		if i > 0 && i%10000 == 0 {
			o.R.log("OnlineMerge.roundSecond 1 %d\n", i)
		}

		currentMachine := o.DeployMap[instance.InstanceId]
		targetMachineId := o.FinalDeployMap[instance.InstanceId]

		//已经部署，直接跳过
		if currentMachine.MachineId == targetMachineId {
			moveAlready++
			continue
		}

		//在临时机器中，本轮不移动
		if o.FreeMachineMap[currentMachine.MachineId] != nil {
			moveFreezing++
			continue
		}

		targetMachine := o.MachineMap[targetMachineId]
		if targetMachine.ConstraintCheck(instance, 1) {
			//迁移到目标机器
			currentMachine.RemoveInstance(instance.InstanceId)
			targetMachine.AddInstance(instance)
			o.DeployMap[instance.InstanceId] = targetMachine

			//留下幽灵
			ghost := instance.CreateGhost()
			currentMachine.AddInstance(ghost)
			ghosts = append(ghosts, ghost)
			ghostsDeploy = append(ghostsDeploy, currentMachine)

			//纪录指令
			o.MoveCommands = append(o.MoveCommands, &InstanceMoveCommand{
				Round:      2,
				InstanceId: instance.InstanceId,
				MachineId:  targetMachine.MachineId,
			})

			moveSuccess++
		} else {
			//迁移尝试失败
			moveRest++
		}
	}

	o.R.log("OnlineMerge.roundSecond 1,machines=%d,deployed=%d,free=%d,"+
		"moveAlready=%d,moveSuccess=%d,moveFreezing=%d,moveRest=%d\n",
		len(o.MachineList), len(o.DeployedMachineList), len(o.FreeMachineList),
		moveAlready, moveSuccess, moveFreezing, moveRest)

	//根据最终状态将必须迁移走的实例移开
	finalMachineMap := make([]*Machine, o.R.MaxMachineId+1)
	for i, m := range o.MachineMap {
		if m != nil {
			finalMachineMap[i] = NewMachine(m.R, m.MachineId, m.Config)
		}
	}
	for instanceId, machineId := range o.FinalDeployMap {
		if machineId > 0 {
			finalMachineMap[machineId].AddInstance(o.InstanceMap[instanceId])
		}
	}

	moveAlready = 0
	moveFreezing = 0
	moveKeep := 0
	moveOther := 0
	moveRest = 0
	for i, instance := range o.InstanceList {
		if i > 0 && i%10000 == 0 {
			o.R.log("OnlineMerge.roundSecond 2 %d\n", i)
		}

		currentMachine := o.DeployMap[instance.InstanceId]
		targetMachineId := o.FinalDeployMap[instance.InstanceId]

		//已经部署，直接跳过
		if currentMachine.MachineId == targetMachineId {
			moveAlready++
			continue
		}

		//在临时机器中，本轮不移动
		if o.FreeMachineMap[currentMachine.MachineId] != nil {
			moveFreezing++
			continue
		}

		//若不可保持位置，迁移走
		finalCurrentMachine := finalMachineMap[currentMachine.MachineId]
		if finalCurrentMachine.ConstraintCheck(instance, 1) {
			//可以保持位置
			finalCurrentMachine.AddInstance(instance)
			moveKeep++
		} else {
			//迁移走
			moved := false
			for _, deployMachine := range o.DeployedMachineList {
				if deployMachine == currentMachine {
					continue
				}

				if !deployMachine.ConstraintCheck(instance, 1) {
					continue
				}

				finalDeployMachine := finalMachineMap[deployMachine.MachineId]
				if !finalDeployMachine.ConstraintCheck(instance, 1) {
					continue
				}

				//迁移实例
				currentMachine.RemoveInstance(instance.InstanceId)
				deployMachine.AddInstance(instance)
				o.DeployMap[instance.InstanceId] = deployMachine

				//更新最终状态
				finalDeployMachine.AddInstance(instance)

				//留下幽灵
				ghost := instance.CreateGhost()
				currentMachine.AddInstance(ghost)
				ghosts = append(ghosts, ghost)
				ghostsDeploy = append(ghostsDeploy, currentMachine)

				//纪录指令
				o.MoveCommands = append(o.MoveCommands, &InstanceMoveCommand{
					Round:      2,
					InstanceId: instance.InstanceId,
					MachineId:  deployMachine.MachineId,
				})

				moveOther++
				moved = true
				break
			}
			if !moved {
				moveRest++
			}
		}
	}

	//删除幽灵实例
	for i, ghost := range ghosts {
		ghostsDeploy[i].RemoveInstance(ghost.InstanceId)
	}

	o.R.log("OnlineMerge.roundSecond 2,machines=%d,deployed=%d,free=%d,"+
		"moveAlready=%d,moveFreezing=%d,moveKeep=%d,moveOther=%d,moveRest=%d\n",
		len(o.MachineList), len(o.DeployedMachineList), len(o.FreeMachineList),
		moveAlready, moveFreezing, moveKeep, moveOther, moveRest)
}

func (o *OnlineMerge) roundFinal() (err error) {
	o.R.log("OnlineMerge.roundFinal start\n")

	//迁移后留下的幽灵实例
	ghosts := make([]*Instance, 0)
	ghostsDeploy := make([]*Machine, 0)

	moveAlready := 0
	moveSuccess := 0
	moveRest := 0
	for i, instance := range o.InstanceList {
		if i > 0 && i%10000 == 0 {
			o.R.log("OnlineMerge.roundFinal %d\n", i)
		}

		currentMachine := o.DeployMap[instance.InstanceId]
		targetMachineId := o.FinalDeployMap[instance.InstanceId]

		//已经部署，直接跳过
		if currentMachine.MachineId == targetMachineId {
			moveAlready++
			continue
		}

		targetMachine := o.MachineMap[targetMachineId]
		if targetMachine.ConstraintCheck(instance, 1) {
			//迁移到目标机器
			currentMachine.RemoveInstance(instance.InstanceId)
			targetMachine.AddInstance(instance)
			o.DeployMap[instance.InstanceId] = targetMachine

			//留下幽灵
			ghost := instance.CreateGhost()
			currentMachine.AddInstance(ghost)
			ghosts = append(ghosts, ghost)
			ghostsDeploy = append(ghostsDeploy, currentMachine)

			//纪录指令
			o.MoveCommands = append(o.MoveCommands, &InstanceMoveCommand{
				Round:      3,
				InstanceId: instance.InstanceId,
				MachineId:  targetMachine.MachineId,
			})

			moveSuccess++
		} else {
			//迁移尝试失败
			moveRest++
		}
	}

	//删除幽灵实例
	for i, ghost := range ghosts {
		ghostsDeploy[i].RemoveInstance(ghost.InstanceId)
	}

	o.R.log("OnlineMerge.roundFinal end,machines=%d,deployed=%d,free=%d,moveAlready=%d,moveSuccess=%d,moveRest=%d\n",
		len(o.MachineList), len(o.DeployedMachineList), len(o.FreeMachineList), moveAlready, moveSuccess, moveRest)

	if moveRest > 0 {
		return fmt.Errorf("OnlineMerge.roundFinal failed,rest=%d\n", moveRest)
	}

	return nil
}

func (o *OnlineMerge) Run() (moveCommands []*InstanceMoveCommand, err error) {
	o.R.log("OnlineMerge.Run\n")

	o.init()

	o.roundFirst()

	o.roundSecond()

	err = o.roundFinal()
	if err != nil {
		return nil, err
	}

	return o.MoveCommands, nil
}
