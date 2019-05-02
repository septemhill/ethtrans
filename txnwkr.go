package main

import (
	"github.com/septemhill/ethacctdb/db"
	"github.com/septemhill/ethacctdb/types"
	"strconv"
	"sync"
	"time"
)

type TxnWorker struct {
	wg        sync.WaitGroup
	next      int64
	done      chan struct{}
	cli       *EtherRPCClient
	workerCnt int
}

func (tw *TxnWorker) getTransactions(blockNumber int64, tc chan<- types.Transaction) bool {
	var b *types.Block
	tw.cli.Request("eth_getBlockByNumber", &b, "0x"+strconv.FormatInt(blockNumber, 16), true)

	if b == nil {
		return false
	}

	//Set timestamp for transaction
	for key := range b.Transactions {
		t, _ := strconv.ParseInt(b.Timestamp, 0, 64)
		b.Transactions[key].Timestamp = time.Unix(t, 0)
		tc <- b.Transactions[key]
	}

	return true
}

func (tw *TxnWorker) Start() {
	worker := func(i int) <-chan types.Transaction {
		tc := make(chan types.Transaction)
		tw.wg.Add(1)

		go func(done <-chan struct{}) {
			defer tw.wg.Done()
			for {
				if ok := tw.getTransactions(int64(i), tc); ok {
					i += tw.workerCnt
				}
				select {
				case <-done:
					close(tc)
					return
				default:
					time.Sleep(time.Millisecond)
				}
			}
		}(tw.done)

		return tc
	}

	merge := func(tcs ...<-chan types.Transaction) <-chan types.Transaction {
		var wg sync.WaitGroup
		out := make(chan types.Transaction, 100)

		wg.Add(len(tcs))

		sender := func(c <-chan types.Transaction) {
			for n := range c {
				out <- n
			}

			wg.Done()
		}

		for _, tc := range tcs {
			go sender(tc)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}

	wc := make([]<-chan types.Transaction, 0)
	for i := 0; i < tw.workerCnt; i++ {
		wc = append(wc, worker(i))
	}

	saveTxns := func(tc <-chan types.Transaction) {
		tw.wg.Add(1)
		defer tw.wg.Done()

		db := db.GetRDBInstance()
		defer db.Close()

		for txn := range tc {
			db.Insert(&txn)
		}
	}

	go saveTxns(merge(wc...))
}

func (tw *TxnWorker) Stop() {
	close(tw.done)
	tw.wg.Wait()
	tw.cli.Close()
}

func NewTxnWorker(startBlock int64, url string, workerCnt int) *TxnWorker {
	tw := &TxnWorker{
		//url:  url,
		workerCnt: workerCnt,
		next:      startBlock,
		done:      make(chan struct{}),
		cli:       NewEtherRPCClient(url),
	}

	return tw
}
