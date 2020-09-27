package main

import (
	"fmt"
	"os"
	"syscall"
)

/**
 * @Struct      Filelock
 * @Description wrapper struct for system filelock object.
 */
type Filelock struct {
	path    string   // filelock absolute path
	filePtr *os.File // filelock file pointer
}

/**
 * @Function    NewFilelock
 * @Description instantiate Filelock struct object.
 * @Param       path - filelock absolute path
 * @Return      Filelock object pointer
 */
func NewFilelock(path string) (*Filelock, error) {
	if FileExist(path) == false {
		tempFilePtr, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("create file failed %s - %s\n", path, err)
		}
		err = tempFilePtr.Close()
		if err != nil {
			return nil, fmt.Errorf("close temp file failed %s - %s\n", path, err)
		}
	}

	return &Filelock{
		path: path,
	}, nil
}

/**
 * @Function    Lock
 * @Description lock the filelock, and block other process require this filelock.
 * @Param
 * @Return      error
 */
func (l *Filelock) Lock() error {
	f, err := os.Open(l.path)
	if err != nil {
		return err
	}
	l.filePtr = f

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		return fmt.Errorf("cannot flock path %s - %s", l.path, err)
	}
	return nil
}

/**
 * @Function    Unlock
 * @Description unlock the filelock.
 * @Param
 * @Return      error
 */
func (l *Filelock) Unlock() error {
	err := syscall.Flock(int(l.filePtr.Fd()), syscall.LOCK_UN)
	if err != nil {
		return fmt.Errorf("flock unlock failed %s - %s", l.path, err)
	}

	err = l.filePtr.Close()
	if err != nil {
		return fmt.Errorf("close filelock file failed %s - %s", l.path, err)
	}
	return nil
}
