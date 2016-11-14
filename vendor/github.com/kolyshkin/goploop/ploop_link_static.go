// +build static_build

package ploop

// #cgo pkg-config: --static ploop
// #cgo LDFLAGS: -static
import "C"
