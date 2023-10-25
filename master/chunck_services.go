package master

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lietoast/catiohead/pb"
	"google.golang.org/grpc"
)

// possible status of a chunck service
const (
	DEAD    = 0
	HEALTHY = 1
	UNKNOWN = 2
	DYING   = 3
)

// config file path
const manifestPath = "./manifest" // test environment

// info of one chunck service
type chunckService struct {
	id              int    // unmodifiable service ID
	address         string // IP address of this service
	port            int
	status          int
	diskFree        uint64           // free disk space remaining in this server
	recentChunckNum int              // chunck created in 60 seconds
	serviceInfoMtx  *sync.RWMutex    // lock for service information
	conn            *grpc.ClientConn // grpc connection to this service
	connMtx         *sync.Mutex      // lock for grpc connection
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

func withChunckServiceInfoLock() chunckServiceOption {
	return func(cs *chunckService) {
		cs.serviceInfoMtx = new(sync.RWMutex)
	}
}

func withChunckServiceConnLock() chunckServiceOption {
	return func(cs *chunckService) {
		cs.connMtx = new(sync.Mutex)
	}
}

func withRecentChunckNum(recentChunckNum int) chunckServiceOption {
	return func(cs *chunckService) {
		cs.recentChunckNum = recentChunckNum
	}
}

func createChunckService(options ...chunckServiceOption) *chunckService {
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

				chunckServices[id] = createChunckService(
					withChunckServiceID(id),
					withChunckServiceAddress(ip),
					withChunckServicePort(port),
					withChunckServiceStatus(UNKNOWN),
					withRecentChunckNum(0),
					withChunckServiceInfoLock(),
					withChunckServiceConnLock(),
				)
			}
		})
}

func init() {
	// init services
	initServicesInfo()
	// check status of services
	checkChunckServicesStatus()
}

func checkChunckServicesStatus() {
	for i := 0; i < len(chunckServices); i++ {
		checkChunckServiceStatus(i)
	}
}

func checkChunckServiceStatus(id int) {
	if service, ok := chunckServices[id]; !ok {
		return
	} else {
		if service.conn == nil {
			service.connMtx.Lock()
			if service.conn == nil {
				conn, err := pb.InsecureConnect(service.address, service.port)
				if err != nil {
					service.conn = nil
					service.setStatus(DEAD)
					service.connMtx.Unlock()
					return
				}
				service.conn = conn
			}
			service.connMtx.Unlock()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		client := pb.NewHealthCheckServiceClient(service.conn)

		resp, err := client.CheckHealthStatus(ctx, &pb.Ping{})
		if err != nil {
			service.connMtx.Lock()
			service.conn = nil
			service.connMtx.Unlock()
			service.setStatus(DEAD)
			return
		}

		service.setDiskFree(resp.DiskFree)
		service.setStatus(int(resp.Status))
		service.setRecentChunckNum(int(resp.RecentChunkNum))
	}
}

func (c *chunckService) setStatus(status int) {
	c.serviceInfoMtx.Lock()
	c.status = status
	c.serviceInfoMtx.Unlock()
}

func (c *chunckService) setDiskFree(diskFree uint64) {
	c.serviceInfoMtx.Lock()
	c.diskFree = diskFree
	c.serviceInfoMtx.Unlock()
}

func (c *chunckService) setRecentChunckNum(recentChunckNum int) {
	c.serviceInfoMtx.Lock()
	c.recentChunckNum = recentChunckNum
	c.serviceInfoMtx.Unlock()
}
