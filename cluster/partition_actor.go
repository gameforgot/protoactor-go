package cluster

import (
	"log"
	"time"

	"github.com/AsynkronIT/gam/actor"
)

var (
	partitionPid = actor.SpawnNamed(actor.FromProducer(newClusterActor()), "partition")
)

func partitionForHost(host string) *actor.PID {
	pid := actor.NewPID(host, "partition")
	return pid
}

func newClusterActor() actor.Producer {
	return func() actor.Actor {
		return &partitionActor{
			partition: make(map[string]*actor.PID),
		}
	}
}

type partitionActor struct {
	partition map[string]*actor.PID
}

func (state *partitionActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("[CLUSTER] Partition started")
	case *ActorPidRequest:
		state.actorPidRequest(msg, context)
	case *clusterStatusJoin:
		state.clusterStatusJoin(msg)
	case *clusterStatusLeave:
		log.Printf("[CLUSTER] Node left %v", msg.node.host)
	case *TakeOwnership:
		state.takeOwnership(msg)
	default:
		log.Printf("[CLUSTER] Partition got unknown message %+v", msg)
	}
}

func (state *partitionActor) actorPidRequest(msg *ActorPidRequest, context actor.Context) {

	pid := state.partition[msg.Name]
	if pid == nil {
		//get a random node
		random := getRandomActivator()

		//send request
		//		log.Printf("[CLUSTER] Telling %v to create %v", random, msg.Name)
		tmp, err := random.RequestFuture(msg, 5*time.Second).Result()
		if err != nil {
			log.Printf("[CLUSTER] Actor PID Request result failed %v on node %v", err, random.Host)
			return
		}
		typed := tmp.(*ActorPidResponse)
		pid = typed.Pid
		state.partition[msg.Name] = pid
	}
	response := &ActorPidResponse{
		Pid: pid,
	}
	context.Respond(response)
}

func (state *partitionActor) clusterStatusJoin(msg *clusterStatusJoin) {
	log.Printf("[CLUSTER] Node joined %v", msg.node.host)
	if list.LocalNode() == nil {
		return
	}

	selfName := list.LocalNode().Name
	for actorID := range state.partition {
		host := getNode(actorID)
		if host != selfName {
			state.transferOwnership(actorID, host)
		}
	}
}

func (state *partitionActor) transferOwnership(actorID string, host string) {
	log.Printf("[CLUSTER] Giving ownership of %v to Node %v", actorID, host)
	pid := state.partition[actorID]
	owner := partitionForHost(host)
	owner.Tell(&TakeOwnership{
		Pid:  pid,
		Name: actorID,
	})
	//we can safely delete this entry as the consisntent hash no longer points to us
	delete(state.partition, actorID)
}

func (state *partitionActor) takeOwnership(msg *TakeOwnership) {
	log.Printf("[CLUSTER] Took ownerhip of %v", msg.Pid)
	state.partition[msg.Name] = msg.Pid
}
