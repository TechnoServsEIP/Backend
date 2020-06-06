// +build darwin dragonfly freebsd netbsd openbsd
<<<<<<< HEAD
=======
// +build !js
>>>>>>> clientGRPCBilling

package logrus

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA

func isTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	return err == nil
}
<<<<<<< HEAD

=======
>>>>>>> clientGRPCBilling
