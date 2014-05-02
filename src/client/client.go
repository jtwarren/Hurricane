package client

import "net/rpc"
import "fmt"
import "sync"
import "time"

const (
	OK          = "OK"
	RESET       = "RESET"
	NO_RESPONSE = "NO_RESPONSE"
)

type Err string

type RegisterArgs struct {
	Me string
}

type RegisterReply struct {
	Id  int64
	Err Err
}

type PingArgs struct {
	Id int64
}

type PingReply struct {
	Err Err
}

type MasterClerk struct {
	mu sync.Mutex

	// (host:port) information
	master string
	me     string
	id     int64
}

func MakeMasterClerk(me string, master string) *MasterClerk {
	ck := new(MasterClerk)

	ck.master = master
	ck.me = me

	return ck
}

func CallRPC(srv string, rpcname string, args interface{}, reply interface{}) bool {
	c, errx := rpc.Dial("tcp", srv)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

// If retry=false, return NO_RESPONSE on failure
func (ck *MasterClerk) Ping(retry bool) Err {
	// id should be at least 1 if registration was successful
	if ck.id <= 0 {
		panic("Cannot ping before registering")
	}

	to := 10 * time.Millisecond
	for {
		args := PingArgs{Id: ck.id}
		reply := PingReply{}
		ok := CallRPC(ck.master, "Master.Ping", &args, &reply)
		if ok {
			// either OK or RESET
			return reply.Err
		}

		if retry == false {
			return NO_RESPONSE
		}

		// exponential backoff
		time.Sleep(to)
		if to < 10*time.Second {
			to *= 2
		}
	}
}

// If retry=false, return 0 upon failure
func (ck *MasterClerk) Register(retry bool) int64 {
	to := 50 * time.Millisecond
	for {
		args := RegisterArgs{Me: ck.me}
		reply := RegisterReply{}
		ok := CallRPC(ck.master, "Master.Register", &args, &reply)

		if ok && reply.Err == OK {
			ck.id = reply.Id
			return reply.Id
		}

		if retry == false {
			return 0
		}

		// exponential backoff
		time.Sleep(to)
		if to < 10*time.Second {
			to *= 2
		}
	}
}

type SegmentInput struct {
	SegmentId      int64
	PartitionIndex int
	WorkerUrl      string
	Index          int
}

type ExecArgs struct {
	Command         string
	Segments        []SegmentInput
	OutputSegmentId int64
	Indices         []int
	Parts           int
}

type ExecReply struct {
	Err Err
}