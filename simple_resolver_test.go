package snowflake_test

import (
	"testing"
	"sync"
	"github.com/NeoGitCrt1/go-snowflake"
)

func TestSimpleResolver(t *testing.T) {
	id, _ := snowflake.AtomicResolver(1)

	if id != 0 {
		t.Error("Sequence should be equal 0")
	}
}

func TestSimpleResolverParallel(t *testing.T) {
	le := 200000
	//snowflake.SetSequenceResolver(snowflake.AtomicResolver)
	ch := make(chan uint16, le)
	var wg sync.WaitGroup
	for i := 0; i < le; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, _ := snowflake.AtomicResolver(1)
			ch <- id
		}()
	}
	wg.Wait()
	close(ch)

	mp := make(map[uint16]bool)
	for id := range ch {
		// if _, ok := mp[id]; ok {
		// 	t.Log(snowflake.ParseID(id))
		// }
		mp[id] = true
	}
	if len(mp) != 1<<snowflake.SequenceLength {
		t.Error("map length should be equal", le, "got:", len(mp))
	}
}

func BenchmarkCombinationSimpleResolverParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = snowflake.AtomicResolver(1)
		}
	})
}

func BenchmarkSimpleResolver(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = snowflake.AtomicResolver(1)
	}
}