package master

import (
	"fmt"
	"sync"
)

// possible status of a chunck service
const (
	DEAD    = 0
	HEALTHY = 1
	UNKNOWN = 2
)

// config file path
const manifestPath = "./manifest" // test environment

// info of one chunck service
type chunckService struct {
	id       int    // unmodifiable service ID
	address  string // IP address of this service
	port     int
	status   int
	diskFree uint64 // free disk space remaining in this server
}

// service ID -> service info
var chunckServices map[int]*chunckService
var servicesInstanceOnce sync.Once

type chunckServiceOption func(*chunckService)

func withChunckServiceID(id int) chunckServiceOption {
	return func(cs *chunckService) {
		cs.id = id
	}
}

func withChunckServiceAddress(address string) chunckServiceOption {
	return func(cs *chunckService) {
		cs.address = address
	}
}

func withChunckServicePort(port int) chunckServiceOption {
	return func(cs *chunckService) {
		cs.port = port
	}
}

func withChunckServiceStatus(status int) chunckServiceOption {
	return func(cs *chunckService) {
		cs.status = status
	}
}

func NewChunckService(options ...chunckServiceOption) *chunckService {
	service := new(chunckService)
	for _, f := range options {
		f(service)
	}
	return service
}

// read chunck service infomation from manifest
func initServicesInfo() {
	servicesInstanceOnce.Do(
		func() {
			services := readManifest(manifestPath)

			for id, address := range services {
				var ip string
				var port int
				n, err := fmt.Sscanf(address, "%s:%d", &ip, &port)
				if n < 2 || err != nil {
					continue
				}

				chunckServices[id] = NewChunckService(
					withChunckServiceID(id),
					withChunckServiceAddress(ip),
					withChunckServicePort(port),
					withChunckServiceStatus(UNKNOWN),
				)
			}
		})
}

func init() {
	// init services
	initServicesInfo()
	// TODO: check status of services
}
