package snowflake

import "sync/atomic"

var lastSeqS uint32
func SimpleResolver(ms int64) (uint16, error) {
	var seq, localSeq uint32

	for {
		localSeq = atomic.LoadUint32(&lastSeqS)
		seq = MaxSequence & (localSeq + 1)
		if atomic.CompareAndSwapUint32(&lastSeqS, localSeq, seq) {
			return uint16(seq), nil
		}
	}
	
	
}