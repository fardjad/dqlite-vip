package learning

import "net"

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func FindFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}
