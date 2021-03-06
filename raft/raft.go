// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package raft

import (
	"errors"

	pb "github.com/pingcap-incubator/tinykv/proto/pkg/eraftpb"
)

// None is a placeholder node ID used when there is no leader.
const None uint64 = 0

// StateType represents the role of a node in a cluster.
type StateType uint64

const (
	StateFollower StateType = iota
	StateCandidate
	StateLeader
)

var stmap = [...]string{
	"StateFollower",
	"StateCandidate",
	"StateLeader",
}

func (st StateType) String() string {
	return stmap[uint64(st)]
}

// ErrProposalDropped is returned when the proposal is ignored by some cases,
// so that the proposer can be notified and fail fast.
var ErrProposalDropped = errors.New("raft proposal dropped")

// Config contains the parameters to start a raft.
type Config struct {
	// ID is the identity of the local raft. ID cannot be 0.
	ID uint64

	// peers contains the IDs of all nodes (including self) in the raft cluster. It
	// should only be set when starting a new raft cluster. Restarting raft from
	// previous configuration will panic if peers is set. peer is private and only
	// used for testing right now.
	peers []uint64

	// ElectionTick is the number of Node.Tick invocations that must pass between
	// elections. That is, if a follower does not receive any message from the
	// leader of current term before ElectionTick has elapsed, it will become
	// candidate and start an election. ElectionTick must be greater than
	// HeartbeatTick. We suggest ElectionTick = 10 * HeartbeatTick to avoid
	// unnecessary leader switching.
	ElectionTick int
	// HeartbeatTick is the number of Node.Tick invocations that must pass between
	// heartbeats. That is, a leader sends heartbeat messages to maintain its
	// leadership every HeartbeatTick ticks.
	HeartbeatTick int

	// Storage is the storage for raft. raft generates entries and states to be
	// stored in storage. raft reads the persisted entries and states out of
	// Storage when it needs. raft reads out the previous state and configuration
	// out of storage when restarting.
	Storage Storage
	// Applied is the last applied index. It should only be set when restarting
	// raft. raft will not return entries to the application smaller or equal to
	// Applied. If Applied is unset when restarting, raft might return previous
	// applied entries. This is a very application dependent configuration.
	Applied uint64
}

func (c *Config) validate() error {
	if c.ID == None {
		return errors.New("cannot use none as id")
	}

	if c.HeartbeatTick <= 0 {
		return errors.New("heartbeat tick must be greater than 0")
	}

	if c.ElectionTick <= c.HeartbeatTick {
		return errors.New("election tick must be greater than heartbeat tick")
	}

	if c.Storage == nil {
		return errors.New("storage cannot be nil")
	}

	return nil
}

// Progress represents a follower’s progress in the view of the leader. Leader maintains
// progresses of all followers, and sends entries to the follower based on its progress.
type Progress struct {
	Match, Next uint64
}

type Raft struct {
	id uint64
	/* persisted state for raft node */
	Term uint64
	Vote uint64
	// the log
	RaftLog *RaftLog
	/* volatile state for leader raft node, reinitialized after election */
	// log replication progress of each peers
	Prs map[uint64]*Progress

	// this peer's role
	State StateType

	// votes records
	votes map[uint64]bool

	// msgs need to send
	msgs []pb.Message

	// the leader id
	Lead uint64

	// heartbeat interval, should send
	heartbeatTimeout int
	// baseline of election interval
	electionTimeout int
	// number of ticks since it reached last heartbeatTimeout.
	// only leader keeps heartbeatElapsed.
	heartbeatElapsed int
	// Ticks since it reached last electionTimeout when it is leader or candidate.
	// Number of ticks since it reached last electionTimeout or received a
	// valid message from current leader when it is a follower.
	electionElapsed int

	// leadTransferee is id of the leader transfer target when its value is not zero.
	// Follow the procedure defined in section 3.10 of Raft phd thesis.
	// (https://web.stanford.edu/~ouster/cgi-bin/papers/OngaroPhD.pdf)
	// (Used in 3A leader transfer)
	leadTransferee uint64

	// Only one conf change may be pending (in the log, but not yet
	// applied) at a time. This is enforced via PendingConfIndex, which
	// is set to a value >= the log index of the latest pending
	// configuration change (if any). Config changes are only allowed to
	// be proposed if the leader's applied index is greater than this
	// value.
	// (Used in 3A conf change)
	PendingConfIndex uint64
}

// newRaft return a raft peer with the given config
func newRaft(c *Config) *Raft {
	if err := c.validate(); err != nil {
		panic(err.Error())
	}
	// Your Code Here (2A).
	hardState, confState, _ := c.Storage.InitialState()
	lo, _ := c.Storage.FirstIndex()
	hi, _ := c.Storage.LastIndex()
	// logs.Infof("lo is %v, hi is %v", lo, hi)
	ents, _ := c.Storage.Entries(lo, hi+1)
	return &Raft{
		id:               c.ID,
		heartbeatTimeout: c.HeartbeatTick,
		electionTimeout:  c.ElectionTick,
		Term:             hardState.Term,
		Vote:             hardState.Vote,
		RaftLog: &RaftLog{
			storage:   c.Storage,
			applied:   c.Applied,
			entries:   ents,
			committed: hardState.Commit,
			stabled:   hardState.Commit,
		},
		votes: make(map[uint64]bool, len(confState.Nodes)),
		Prs:   make(map[uint64]*Progress, len(confState.Nodes)),
	}
}

// sendAppend sends an append RPC with new entries (if any) and the
// current commit index to the given peer. Returns true if a message was sent.
func (r *Raft) sendAppend(to uint64) bool {
	// Your Code Here (2A).
	return false
}

// sendHeartbeat sends a heartbeat RPC to the given peer.
func (r *Raft) sendHeartbeat(to uint64) {
	// Your Code Here (2A).
	// heartbeatTyp := pb.MessageType_MsgHeartbeat
	// if r.State != StateLeader {
	// 	heartbeatTyp = pb.MessageType_MsgHeartbeatResponse
	// }
	r.msgs = append(r.msgs, pb.Message{
		MsgType: pb.MessageType_MsgHeartbeat,
		To:      to,
		From:    r.id,
		Term:    r.Term,
	})
}

func (r *Raft) sendRequestVote(to uint64) {
	lastLogIndex := r.RaftLog.LastIndex()
	lastLogTerm, _ := r.RaftLog.Term(lastLogIndex)
	for peerID := range r.Prs {
		r.msgs = append(r.msgs, pb.Message{
			MsgType: pb.MessageType_MsgRequestVote,
			To:      peerID,
			From:    r.id,
			Term:    r.Term,
			Index:   lastLogIndex,
			LogTerm: lastLogTerm,
		})
	}
}

func (r *Raft) sendHeartbeatResp(to uint64, rej bool) {
	r.msgs = append(r.msgs, pb.Message{
		MsgType: pb.MessageType_MsgHeartbeatResponse,
		To:      to,
		From:    r.id,
		Term:    r.Term,
		Reject:  rej,
	})
}

func (r *Raft) sendRequestVoteResp(votedFor uint64, rej bool) {
	r.msgs = append(r.msgs, pb.Message{
		MsgType: pb.MessageType_MsgRequestVoteResponse,
		To:      votedFor,
		From:    r.id,
		Term:    r.Term,
		Reject:  rej,
	})
	// votes[votedFor] = true
}

// tick advances the internal logical clock by a single tick.
func (r *Raft) tick() {
	// Your Code Here (2A).
	r.Term++
	r.Vote = None
	r.votes = make(map[uint64]bool)
}

// updateTerm updates current term to given term only when
// cur term <= given term, and reset the vote
func (r *Raft) updateTerm(to uint64) {
	if to < r.Term {
		return
	}
	r.Term = to
	// clear votedFor & votes track map
	r.Vote = None
	r.votes = make(map[uint64]bool)
}

// becomeFollower transform this peer's state to Follower
func (r *Raft) becomeFollower(term uint64, lead uint64) {
	// Your Code Here (2A).
	r.State = StateFollower
	r.Lead = lead
	// r.Term = term
	r.updateTerm(term)
}

// becomeCandidate transform this peer's state to candidate
func (r *Raft) becomeCandidate() {
	// Your Code Here (2A).
	r.State = StateCandidate
	r.Lead = None
}

// becomeLeader transform this peer's state to leader
func (r *Raft) becomeLeader() {
	// Your Code Here (2A).
	// NOTE: Leader should propose a noop entry on its term
}

// Step the entrance of handle message, see `MessageType`
// on `eraftpb.proto` for what msgs should be handled
func (r *Raft) Step(m pb.Message) error {
	// Your Code Here (2A).
	switch r.State {
	case StateFollower:
		if m.MsgType == pb.MessageType_MsgHeartbeat {
			r.handleHeartbeat(m)
		}
	case StateCandidate:
	case StateLeader:
		if m.MsgType == pb.MessageType_MsgHeartbeatResponse {
			r.handleHearbeatResp(m)
		}
	}
	return nil
}

func (r *Raft) handleRequestVote(m pb.Message) {
	if r.Term > m.Term {
		r.sendRequestVoteResp(None, true)
		return
	}
	if r.Term < m.Term {
		r.becomeFollower(m.Term, None)
	}
	lastLogIndex := r.RaftLog.LastIndex()
	lastLogTerm, err := r.RaftLog.Term(lastLogIndex)
	if err != nil {
		panic("invalid lastLogIndex, raft state machine is in chaos state")
	}
	if (r.Vote == None || r.Vote == m.From) &&
		(lastLogIndex <= m.Index && lastLogTerm <= m.LogTerm) {
		r.sendRequestVoteResp(m.From, false)
		return
	}
	// reject the out-of-date RequestVote
	r.sendRequestVoteResp(None, true)
}

func (r *Raft) handleRequestVoteResp(m pb.Message) {

}

// handleAppendEntries handle AppendEntries RPC request
func (r *Raft) handleAppendEntries(m pb.Message) {
	// Your Code Here (2A).
}

// handleHeartbeat handle Heartbeat RPC request
func (r *Raft) handleHeartbeat(m pb.Message) {
	// Your Code Here (2A).
	// if Term of rcvd req > cur raft Term, override it
	// with rcvd Term. Discard the req when Term of
	// req < cur raft Term
	if r.Term < m.Term {
		r.updateTerm(m.Term)
		r.becomeFollower(m.Term, m.From)
		r.sendHeartbeatResp(m.From, false)
		return
	} else if r.Term == m.Term {
		r.sendHeartbeatResp(m.From, false)
		return
	}
	r.sendHeartbeatResp(m.From, true)
}

func (r *Raft) handleHearbeatResp(m pb.Message) {
	// when leader rcvd higher term heartbeat response
	if r.Term < m.Term {
		r.updateTerm(m.Term)
		r.becomeFollower(m.Term, m.From)
		return
	}
}

func (r *Raft) handleAppendEntriesResp(m pb.Message) {

}

// handleSnapshot handle Snapshot RPC request
func (r *Raft) handleSnapshot(m pb.Message) {
	// Your Code Here (2C).
}

// addNode add a new node to raft group
func (r *Raft) addNode(id uint64) {
	// Your Code Here (3A).
}

// removeNode remove a node from raft group
func (r *Raft) removeNode(id uint64) {
	// Your Code Here (3A).
}
