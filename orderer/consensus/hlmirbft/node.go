/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hlmirbft

import (
	"bytes"
	"crypto"
	"sync"

	"github.com/fly2plan/fabric-protos-go/orderer/hlmirbft"
	"github.com/hyperledger-labs/mirbft"
	"github.com/hyperledger-labs/mirbft/pkg/eventlog"
	"github.com/hyperledger-labs/mirbft/pkg/pb/msgs"
	"github.com/hyperledger-labs/mirbft/pkg/reqstore"
	"github.com/hyperledger-labs/mirbft/pkg/simplewal"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/protoutil"

	"code.cloudfoundry.org/clock"
	"github.com/hyperledger/fabric/common/flogging"
	"go.etcd.io/etcd/raft/raftpb"
)

type node struct {
	chainID string
	logger  *flogging.FabricLogger
	metrics *Metrics

	unreachableLock sync.RWMutex
	unreachable     map[uint64]struct{}

	config      *mirbft.Config
	WALDir      string
	ReqStoreDir string

	rpc RPC

	chain *Chain

	clock clock.Clock

	metadata *hlmirbft.BlockMetadata

	mirbft.Node
}

func (n *node) start(fresh, join bool) {
	n.logger.Debugf("Starting mirbft node: #peers: %v", len(n.metadata.ConsenterIds))
	if fresh {
		if join {
			n.logger.Info("Starting mirbft node to join an existing channel")
		} else {
			n.logger.Info("Starting mirbft node as part of a new channel")
		}

		// Checking if the configuration settings have been passed correctly.
		wal, err := simplewal.Open(n.WALDir)
		if err != nil {
			n.logger.Error(err, "Failed to create WAL")
		}
		reqStore, err := reqstore.Open(n.ReqStoreDir)
		if err != nil {
			n.logger.Error(err, "Failed to create request store")
		}
		node, err := mirbft.NewNode(
			n.chain.MirBFTID,
			n.config,
			&mirbft.ProcessorConfig{
				Link:         n,
				Hasher:       crypto.SHA256,
				App:          n.chain,
				WAL:          wal,
				RequestStore: reqStore,
				Interceptor:  eventlog.NewRecorder(n.chain.MirBFTID, &bytes.Buffer{}),
			},
		)
		if err != nil {
			n.logger.Error(err, "Failed to create mirbft node")
		} else {
			n.Node = *node
		}

		initialNetworkState := InitialNetworkState(n.metadata.ConsenterIds)
		// TODO(harrymknight) Tick interval is fixed. Perhaps introduce TickInterval field in configuration options
		err = n.ProcessAsNewNode(n.chain.doneC, n.clock.NewTicker(10).C(), initialNetworkState, []byte("first"))

	} else {
		n.logger.Info("Restarting mirbft node")
		n.RestartProcessing(n.chain.doneC, n.clock.NewTicker(10).C())
	}
}

func InitialNetworkState(ConsenterIds []uint64) *msgs.NetworkState {
	var nodes []uint64
	for i := 0; i < len(ConsenterIds); i++ {
		nodes[i] = ConsenterIds[i]
	}

	numberOfBuckets := int32(len(ConsenterIds))
	checkpointInterval := numberOfBuckets * 5
	maxEpochLength := checkpointInterval * 10

	// TODO(harrymknight) The width of a client window is fixed.
	//  Could optimise by varying according to the checkpoint interval and batch size
	clients := make([]*msgs.NetworkState_Client, len(ConsenterIds))
	for i := range clients {
		clients[i] = &msgs.NetworkState_Client{
			Id:           ConsenterIds[i],
			Width:        100,
			LowWatermark: 0,
		}
	}

	return &msgs.NetworkState{
		Config: &msgs.NetworkState_Config{
			Nodes:              nodes,
			F:                  int32((len(ConsenterIds) - 1) / 3),
			NumberOfBuckets:    numberOfBuckets,
			CheckpointInterval: checkpointInterval,
			MaxEpochLength:     uint64(maxEpochLength),
		},
		Clients: clients,
	}
}

// TODO(harry_knight) The logic contained in the infinite for loops should be retained.
// 	It serves to start, manage, and respond to the internal clock of the FSM.
// 	Auxiliary calls should be adapted to occur during block genesis/orderer service startup.
func (n *node) run(campaign bool) {
	/*electionTimeout := n.tickInterval.Seconds() * float64(n.config.ElectionTick)
	halfElectionTimeout := electionTimeout / 2

	raftTicker := n.clock.NewTicker(n.tickInterval)

	if s := n.storage.Snapshot(); !raft.IsEmptySnap(s) {
		n.chain.snapC <- &s
	}

	elected := make(chan struct{})
	if campaign {
		n.logger.Infof("This node is picked to start campaign")
		go func() {
			// Attempt campaign every two HeartbeatTimeout elapses, until leader is present - either this
			// node successfully claims leadership, or another leader already existed when this node starts.
			// We could do this more lazily and exit proactive campaign once transitioned to Candidate state
			// (not PreCandidate because other nodes might not have started yet, in which case PreVote
			// messages are dropped at recipients). But there is no obvious reason (for now) to be lazy.
			//
			// 2*HeartbeatTick is used to avoid excessive campaign when network latency is significant and
			// Raft term keeps advancing in this extreme case.
			campaignTicker := n.clock.NewTicker(n.tickInterval * time.Duration(n.config.HeartbeatTick) * 2)
			defer campaignTicker.Stop()

			for {
				select {
				case <-campaignTicker.C():
					n.Campaign(context.TODO())
				case <-elected:
					return
				case <-n.chain.doneC:
					return
				}
			}
		}()
	}

	var notifyLeaderChangeC chan uint64

	for {
		select {
		case <-raftTicker.C():
			// grab raft Status before ticking it, so `RecentActive` attributes
			// are not reset yet.
			status := n.Status()

			n.Tick()
			n.tracker.Check(&status)

		case rd := <-n.Ready():
			startStoring := n.clock.Now()
			if err := n.storage.Store(rd.Entries, rd.HardState, rd.Snapshot); err != nil {
				n.logger.Panicf("Failed to persist etcd/raft data: %s", err)
			}
			duration := n.clock.Since(startStoring).Seconds()
			n.metrics.DataPersistDuration.Observe(float64(duration))
			if duration > halfElectionTimeout {
				n.logger.Warningf("WAL sync took %v seconds and the network is configured to start elections after %v seconds. Your disk is too slow and may cause loss of quorum and trigger leadership election.", duration, electionTimeout)
			}

			if !raft.IsEmptySnap(rd.Snapshot) {
				n.chain.snapC <- &rd.Snapshot
			}

			if notifyLeaderChangeC != nil && rd.SoftState != nil {
				if l := atomic.LoadUint64(&rd.SoftState.Lead); l != raft.None {
					select {
					case notifyLeaderChangeC <- l:
					default:
					}

					notifyLeaderChangeC = nil
				}
			}

			// skip empty apply
			if len(rd.CommittedEntries) != 0 || rd.SoftState != nil {
				n.chain.applyC <- apply{rd.CommittedEntries, rd.SoftState}
			}

			if campaign && rd.SoftState != nil {
				leader := atomic.LoadUint64(&rd.SoftState.Lead) // etcdraft requires atomic access to this var
				if leader != raft.None {
					n.logger.Infof("Leader %d is present, quit campaign", leader)
					campaign = false
					close(elected)
				}
			}

			n.Advance()

			// TODO(jay_guo) leader can write to disk in parallel with replicating
			// to the followers and them writing to their disks. Check 10.2.1 in thesis
			n.send(rd.Messages)

		case notifyLeaderChangeC = <-n.subscriberC:

		case <-n.chain.haltC:
			raftTicker.Stop()
			n.Stop()
			n.storage.Close()
			n.logger.Infof("Raft node stopped")
			close(n.chain.doneC) // close after all the artifacts are closed
			return
		}
	}*/
}

func (n *node) Send(dest uint64, msg *msgs.Msg) {
	n.unreachableLock.RLock()
	defer n.unreachableLock.RUnlock()

	msgBytes := protoutil.MarshalOrPanic(msg)
	err := n.rpc.SendConsensus(dest, &orderer.ConsensusRequest{Channel: n.chainID, Payload: msgBytes})
	if err != nil {
		n.logSendFailure(dest, err)
	} else if _, ok := n.unreachable[dest]; ok {
		n.logger.Infof("Successfully sent StepRequest to %d after failed attempt(s)", dest)
		delete(n.unreachable, dest)
	}
}

// If this is called on leader, it picks a node that is
// recently active, and attempt to transfer leadership to it.
// If this is called on follower, it simply waits for a
// leader change till timeout (ElectionTimeout).
func (n *node) abdicateLeader(currentLead uint64) {
	/*status := n.Status()

	if status.Lead != raft.None && status.Lead != currentLead {
		n.logger.Warn("Leader has changed since asked to transfer leadership")
		return
	}

	// register a leader subscriberC
	notifyc := make(chan uint64, 1)
	select {
	case n.subscriberC <- notifyc:
	case <-n.chain.doneC:
		return
	}

	// Leader initiates leader transfer
	if status.RaftState == raft.StateLeader {
		var transferee uint64
		for id, pr := range status.Progress {
			if id == status.ID {
				continue // skip self
			}

			if pr.RecentActive && !pr.Paused {
				transferee = id
				break
			}

			n.logger.Debugf("Node %d is not qualified as transferee because it's either paused or not active", id)
		}

		if transferee == raft.None {
			n.logger.Errorf("No follower is qualified as transferee, abort leader transfer")
			return
		}

		n.logger.Infof("Transferring leadership to %d", transferee)
		n.TransferLeadership(context.TODO(), status.ID, transferee)
	}

	timer := n.clock.NewTimer(time.Duration(n.config.ElectionTick) * n.tickInterval)
	defer timer.Stop() // prevent timer leak

	select {
	case <-timer.C():
		n.logger.Warn("Leader transfer timeout")
	case l := <-notifyc:
		n.logger.Infof("Leader has been transferred from %d to %d", currentLead, l)
	case <-n.chain.doneC:
	}*/
}

func (n *node) logSendFailure(dest uint64, err error) {
	if _, ok := n.unreachable[dest]; ok {
		n.logger.Debugf("Failed to send StepRequest to %d, because: %s", dest, err)
		return
	}

	n.logger.Errorf("Failed to send StepRequest to %d, because: %s", dest, err)
	n.unreachable[dest] = struct{}{}
}

func (n *node) takeSnapshot(index uint64, cs raftpb.ConfState, data []byte) {
	/*if err := n.storage.TakeSnapshot(index, cs, data); err != nil {
		n.logger.Errorf("Failed to create snapshot at index %d: %s", index, err)
	}*/
}

func (n *node) lastIndex() uint64 {
	/*i, _ := n.storage.ram.LastIndex()
	return i*/
	return 0
}
