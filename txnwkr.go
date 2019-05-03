package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/septemhill/ethacctdb/db"
	"github.com/septemhill/ethacctdb/types"
)

type TxnWorker struct {
	wg         sync.WaitGroup
	done       chan struct{}
	cli        *EtherRPCClient
	workerCnt  int
	lastBlocks []int64
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
	worker := func(i int64) <-chan types.Transaction {
		tc := make(chan types.Transaction)
		tw.wg.Add(1)

		go func(done <-chan struct{}) {
			defer tw.wg.Done()
			for {
				if ok := tw.getTransactions(int64(i), tc); ok {
					i += int64(tw.workerCnt)
				}
				select {
				case <-done:
					tw.lastBlocks = append(tw.lastBlocks, int64(i))
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
		wc = append(wc, worker(tw.lastBlocks[i]))
		//wc = append(wc, worker(i))
	}

	tw.lastBlocks = make([]int64, 0)

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

func (tw *TxnWorker) processedRecord() {
	tpr := &types.TransactionProcessRecord{
		Workers:    tw.workerCnt,
		LastBlocks: tw.lastBlocks,
	}

	buff := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buff)
	enc.Encode(tpr)

	ioutil.WriteFile("processedData.json", buff.Bytes(), 0664)
}

func (tw *TxnWorker) loadProcessedRecord() bool {
	fd, err := os.OpenFile("processedData.json", syscall.O_RDONLY, 0664)

	if err != nil {
		return false
	}

	defer fd.Close()

	b, _ := ioutil.ReadAll(fd)
	tpr := types.TransactionProcessRecord{}
	buff := bytes.NewBuffer(b)
	dec := json.NewDecoder(buff)
	dec.Decode(&tpr)

	tw.workerCnt = tpr.Workers
	tw.lastBlocks = tpr.LastBlocks

	return false
}

func (tw *TxnWorker) Stop() {
	close(tw.done)
	tw.wg.Wait()
	tw.processedRecord()
	tw.cli.Close()
}

func createTxnTable() {
	db := db.GetRDBInstance()

	if err := db.CreateTable(&types.Transaction{}, &orm.CreateTableOptions{}); err == nil {
		fmt.Println("create create successful")
	}
}

func NewTxnWorker(startBlock int64, url string, workerCnt int) *TxnWorker {
	createTxnTable()

	tw := &TxnWorker{
		//url:  url,
		workerCnt: workerCnt,
		done:      make(chan struct{}),
		cli:       NewEtherRPCClient(url),
	}

	if !tw.loadProcessedRecord() {
		for i := 0; i < tw.workerCnt; i++ {
			tw.lastBlocks = append(tw.lastBlocks, int64(i))
		}
	}

	return tw
}
