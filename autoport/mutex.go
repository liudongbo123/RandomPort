package main

import (
	"time"
)

/**
 * @Struct      PortMutex
 * @Description port mutex struct.
 */
type PortMutex struct {
	DelaySecond int64
	PortMap     *map[int]int64
}

/**
 * @Function    NewPortMutex
 * @Description create port mutex var inst.
 * @Param       delaySecond - time by second to hold up the tcp port
 * @Return      PortMutex object pointer
 */
func NewPortMutex(delaySecond int64) (*PortMutex, error) {
	tempMap := make(map[int]int64)
	return &PortMutex{
		DelaySecond: delaySecond,
		PortMap:     &tempMap,
	}, nil
}

/**
 * @Function    Update
 * @Description update port mutex objet.
 * @Param       port - linux port number
 *              utt  - unavailable To this Time, 0 by default
 * @Return      error
 */
func (p *PortMutex) Update(port int, utt int64) error {
	portMap := *(p.PortMap)
	nowTS := time.Now().Unix()

	oldUtt, found := portMap[port]
	if found {
		newUtt := Int64Max(utt, oldUtt)
		if nowTS < newUtt {
			portMap[port] = newUtt
		} else {
			delete(portMap, port)
		}
	} else {
		if nowTS < utt {
			portMap[port] = utt
		}
	}
	return nil
}

/**
 * @Function    Filter
 * @Description update port mutex objet.
 * @Param       portList - port list waiting to be filtered
 * @Return      *[]int - available port list
 */
func (p *PortMutex) Filter(portList []int) []int {
	result := make([]int, 0, len(portList))

	if portList == nil || len(portList) == 0 {
		return result
	}

	for _, port := range portList {
		_, found := (*(p.PortMap))[port]
		if found {
			continue
		} else {
			result = append(result, port)
		}
	}
	return result
}
