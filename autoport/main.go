package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const DelaySecond = 120    // set port unavailable until 120 seconds later.
const MaxPortRetrieve = 40 // get available port count from linux once.
const NeedPortCount = 2    // port count we need.
const MaxRetryTimes = 30   // max retry times.

const FileLock = "/tmp/filelock/lock" // filelock file
const PortFile = "/tmp/filelock/port" // port record for mutex

const AppPort = "/tmp/kuaiyun/appPort"     // app port
const DebugPort = "/tmp/kuaiyun/debugPort" // debug port

/**
 * 通过文件锁控制并发
 */
func (p *PortMutex) FilterByFilelock(ports []int) (int, int, error) {
	// filelock file
	if FileExist(PortFile) == false {
		tempFilePtr, _ := os.Create(PortFile)
		err := tempFilePtr.Close()
		if err != nil {
			return 0, 0, fmt.Errorf("create port file failed %s - %s\n", PortFile, err)
		}
	}

	// get lock
	lock, _ := NewFilelock(FileLock)
	for i := 0; i < MaxRetryTimes; i++ {
		err := lock.Lock()
		if err != nil {
			fmt.Printf("get filelock failed, retry in 10ms.\n")
			time.Sleep(time.Millisecond * 10)
		} else {
			break
		}

		if i == MaxRetryTimes {
			return 0, 0, fmt.Errorf("get filelock failed after %d retry.\n", MaxRetryTimes)
		}
	}

	// critical section begin
	portFilePtr, _ := os.Open(PortFile)
	reader := bufio.NewReader(portFilePtr)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		row := string(line[:])
		fields := strings.Split(row, ":")
		if fields == nil || len(fields) != 2 {
			fmt.Printf("skip wrong format line : %s\n", row)
			continue
		}

		port, _ := strconv.Atoi(fields[0])
		utt, _ := strconv.ParseInt(fields[1], 10, 64)
		_ = p.Update(port, utt)
	}

	for port := range ports {
		_ = p.Update(port, 0)
	}

	var result []int
	for _, port := range ports {
		if _, found := (*(p.PortMap))[port]; found {
			continue
		} else {
			result = append(result, port)
			if len(result) == NeedPortCount {
				break
			}
		}
	}

	// get enough ports
	if len(result) == NeedPortCount {
		(*(p.PortMap))[result[0]] = time.Now().Unix() + DelaySecond
		(*(p.PortMap))[result[1]] = time.Now().Unix() + DelaySecond
	} else {
		_ = portFilePtr.Close()
		return 0, 0, fmt.Errorf("there is not enough ports.\n")
	}
	_ = portFilePtr.Close()

	// write back
	portFilePtr, _ = os.OpenFile(PortFile, os.O_WRONLY, 0)
	_ = portFilePtr.Truncate(0)
	for port, utt := range *(p.PortMap) {
		p := strconv.Itoa(port)
		u := strconv.FormatInt(utt, 10)
		_, _ = portFilePtr.WriteString(strings.Join([]string{p, u}, ":") + "\n")
	}
	_ = portFilePtr.Close()

	// critical section end
	_ = lock.Unlock()
	return result[0], result[1], nil
}

func main() {
	portMutex, _ := NewPortMutex(DelaySecond)

	// update unavailable port map
	ports, _ := GetFreeTcpPorts(MaxPortRetrieve)
	if ports == nil || len(ports) < NeedPortCount {
		panic("there is not enough available port in linux.")
	}
	appPort, debugPort, _ := portMutex.FilterByFilelock(ports)

	// write port to volume
	app, _ := os.Create(AppPort)
	debug, _ := os.Create(DebugPort)

	_, _ = io.WriteString(app, strconv.Itoa(appPort))
	_, _ = io.WriteString(debug, strconv.Itoa(debugPort))

	_ = app.Close()
	_ = debug.Close()
	fmt.Printf("retrieved avaliable ports: %d %d\n", appPort, debugPort)
}
