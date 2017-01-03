// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package main

func getRealPath(path string) string {
	return path
}
