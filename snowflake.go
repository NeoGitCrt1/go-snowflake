// Package snowflake is a network service for generating unique ID numbers at high scale with some simple guarantees.
// The first bit is unused sign bit.
// The second part consists of a 41-bit timestamp (milliseconds) whose value is the offset of the current time relative to a certain time.
// The 5 bits of the third and fourth parts represent data center and worker, and max value is 2^5 -1 = 31.
// The last part consists of 12 bits, its means the length of the serial number generated per millisecond per working node, a maximum of 2^12 -1 = 4095 IDs can be generated in the same millisecond.
// In a distributed environment, five-bit datacenter and worker mean that can deploy 31 datacenters, and each datacenter can deploy up to 31 nodes.
// The binary length of 41 bits is at most 2^41 -1 millisecond = 69 years. So the snowflake algorithm can be used for up to 69 years, In order to maximize the use of the algorithm, you should specify a start time for it.
package snowflake

import (
	"time"
)

// These constants are the bit lengths of snowflake ID parts.
const (
	TimestampLength = 41
	MachineIDLength = 10
	SequenceLength  = 12
	MaxSequence     = 1<<SequenceLength - 1
	MaxTimestamp    = 1<<TimestampLength - 1
	MaxMachineID    = 1<<MachineIDLength - 1

	machineIDMoveLength = SequenceLength
	timestampMoveLength = MachineIDLength + SequenceLength
)

// SequenceResolver the snowflake sequence resolver.
//
// When you want use the snowflake algorithm to generate unique ID, You must ensure: The sequence-number generated in the same millisecond of the same node is unique.
// Based on this, we create this interface provide following reslover:
//   AtomicResolver : base sync/atomic (by default).
type SequenceResolver func(ms int64) (uint16, error)

// default start time is 2008-11-10 23:00:00 UTC, why ? In the playground the time begins at 2009-11-10 23:00:00 UTC.
// It's can run on golang playground.
// default machineID is 0
var (
	machineID   = 0
	stTimestamp = currentMillis(time.Date(2008, 11, 10, 23, 0, 0, 0, time.UTC))
)

// below for reference
// go func() {
// 	for {
// 		value := rand.Intn(MaxRandomNumber)

// 		select {
// 		case <- stopCh:
// 			return
// 		case dataCh <- value:
// 		}
// 	}
// }()
var idChan chan uint64 = make(chan uint64, 200)

func init() {
	go func() {
		mId := &machineID
		st := &stTimestamp
		var now int64 = currentMillis(time.Now())
		var seq int = 0
		var df int
		var pre int
		for {
			if seq == 0 {
				now++
				df = int(now - *st)
				pre = (df << timestampMoveLength) | (*mId << machineIDMoveLength)
			}
			seq = MaxSequence & (seq + 1)
			idChan <- uint64(pre | seq)
		}
	}()

}

// ID use ID to generate snowflake id .
// This function is thread safe.
func ID() uint64 {
	return <-idChan
}

// SetStartTime set the start time for snowflake algorithm.
//
// It will panic when:
//   s IsZero
//   s > current millisecond
//   current millisecond - s > 2^41(69 years).
// This function is thread-unsafe, recommended you call him in the main function.
func SetStartTime(s time.Time) {
	s = s.UTC()

	if s.IsZero() {
		panic("The start time cannot be a zero value")
	}

	if s.After(time.Now()) {
		panic("The s cannot be greater than the current millisecond")
	}

	// Because s must after now, so the `df` not < 0.
	df := currentMillis(time.Now()) - currentMillis(s)
	if df > MaxTimestamp {
		panic("The maximum life cycle of the snowflake algorithm is 69 years")
	}
	stTimestamp = currentMillis(s)
	var v uint64 = ParseID(<-idChan).Timestamp
	for s := range idChan {
		t := ParseID(s).Timestamp
		if t != v {
			break
		}
		v = t
	}
}

func SetDateCenterIdWorkerId(d uint8, w uint8) {
	SetMachineID(conbineDandW(d, w, 5))
}
func SetDateCenterIdWorkerIdWithLen(d uint8, w uint8, wl uint8) {
	SetMachineID(conbineDandW(d, w, wl))
}

func conbineDandW(d uint8, w uint8, wl uint8) uint16 {
	if ( wl > MachineIDLength) {
		panic("wl should less than 10")
	}
	return (uint16(d)<< (wl)) | uint16(w)
}

// SetMachineID specify the machine ID. It will panic when machineid > max limit for 2^10-1.
// This function is thread-unsafe, recommended you call him in the main function.
func SetMachineID(m uint16) {
	if m > MaxMachineID {
		panic("The machineid cannot be greater than 1023")
	}
	machineID = int(m)
	var v uint64 = ParseID(<-idChan).MachineID
	for s := range idChan {
		t := ParseID(s).MachineID
		if t != v {
			break
		}
		v = t
	}
}

// SID snowflake id
type SID struct {
	Sequence  uint64
	MachineID uint64
	Timestamp uint64
	ID        uint64
}

// GenerateTime snowflake generate at, return a UTC time.
func (id *SID) GenerateTime() time.Time {
	ms := stTimestamp + int64(id.Timestamp)

	return time.Unix(0, (ms * int64(time.Millisecond))).UTC()
}

// ParseID parse snowflake it to SID struct.
func ParseID(id uint64) SID {
	time := id >> (SequenceLength + MachineIDLength)
	sequence := id & MaxSequence
	machineID := (id & (MaxMachineID << SequenceLength)) >> SequenceLength

	return SID{
		ID:        id,
		Sequence:  sequence,
		MachineID: machineID,
		Timestamp: time,
	}
}

// currentMillis get current millisecond.
func currentMillis(t time.Time) int64 {
	return t.UTC().UnixNano() / 1e6
}
