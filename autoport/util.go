package main

import (
	"fmt"
	"net"
	"os"
)

/**
 * @Function    FileExist
 * @Description check the file specified by path.
 * @Param       path - filename(absolute path)
 * @Return      true  - file exists
 *              false - file doesn't exist
 */
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

/**
 * @Function    Int64Max
 * @Description compare and return the bigger.
 * @Param       operatorLeft  - first int64 var
 *              operatorRight - second int64 var
 * @Return      the bigger
 */
func Int64Max(operatorLeft, operatorRight int64) int64 {
	if operatorLeft > operatorRight {
		return operatorLeft
	} else {
		return operatorRight
	}
}

/**
 * @Function    GetFreePorts
 * @Description apply available port from System.
 * @Param       count - port number we need
 * @Return      []int - port slice
 *              error
 */
func GetFreeTcpPorts(count int) ([]int, error) {
	portSlice := make([]int, 0, count)
	listenerSlice := make([]*net.TCPListener, 0, count)

	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return portSlice, err
		}
		listenerSlice = append(listenerSlice, listener)
		portSlice = append(portSlice, listener.Addr().(*net.TCPAddr).Port)
	}

	for _, listener := range listenerSlice {
		err := listener.Close()
		return portSlice, fmt.Errorf("close tcp listener failed. %s\n", err)
	}
	return portSlice, nil
}
