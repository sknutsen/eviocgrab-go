package eviocgrabgo

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"golang.org/x/sys/unix"
)

const (
	_IOC_DIRSHIFT  = 30
	_IOC_TYPESHIFT = 8
	_IOC_NRSHIFT   = 0
	_IOC_SIZESHIFT = 16

	_IOC_NONE  = 0
	_IOC_WRITE = 1

	_UNGRAB = 0
	_GRAB   = 1
)

type IOCTL struct {
	_EVIOCGRAB uint
}

func Init() (*IOCTL, error) {
	var ioc *IOCTL = &IOCTL{}

	path, err := exec.LookPath("grep")
	if err != nil {
		fmt.Printf("didn't find 'grep' executable\n")
	} else {
		fmt.Printf("'grep' executable is in '%s'\n", path)
	}

	cmd := exec.Command("grep", "-rn", "EVIOCGRAB", "/usr/include/linux/")

	println(cmd.String())

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("RUN - Error: %s\n", err)
		return nil, err
	}

	re := regexp.MustCompile(`(?m)(0x\d*)`)
	match := re.FindString(string(stdout))
	u, err := strconv.ParseUint(match, 16, 64)
	if err != nil {
		return nil, err
	}

	ioc._EVIOCGRAB = uint(u)

	return ioc, nil
}

func (i *IOCTL) EVIOCGRAB() uint {
	return _IOC(_IOC_WRITE, 'E', i._EVIOCGRAB, 0)
}

func _IOC(dir, typ, nr, size uint) uint {
	return ((dir << _IOC_DIRSHIFT) |
		(typ << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT) |
		(size << _IOC_SIZESHIFT))
}

func (i *IOCTL) Grab(fd *os.File) error {
	return unix.IoctlSetInt(int(fd.Fd()), i.EVIOCGRAB(), _GRAB)
}

func (i *IOCTL) UnGrab(fd *os.File) error {
	return unix.IoctlSetInt(int(fd.Fd()), i.EVIOCGRAB(), _UNGRAB)
}
