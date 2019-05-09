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

func (tw *TxnWorker) getReceipts(txnHash string, rc chan<- types.Receipt) bool {
	var r *types.Receipt
	tw.cli.Request("eth_getTransactionReceipt", &r, txnHash)

	if r == nil {
		return false
	}

	rc <- *r

	return true
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

func (tw *TxnWorker) getTransactions2(blockNumber int64, tc chan<- types.Transaction, rc chan<- types.Receipt) bool {
	var b *types.Block
	tw.cli.Request("eth_getBlockByNumber", &b, "0x"+strconv.FormatInt(blockNumber, 16), true)

	if b == nil {
		return false
	}

	//Set timestamp for transaction
	for key := range b.Transactions {
		t, _ := strconv.ParseInt(b.Timestamp, 0, 64)
		b.Transactions[key].Timestamp = time.Unix(t, 0)
		tw.getReceipts(b.Transactions[key].Hash, rc)
		tc <- b.Transactions[key]
	}

	return true
}

func (tw *TxnWorker) Start() {
	worker := func(i int64) (<-chan types.Transaction, <-chan types.Receipt) {
		tc := make(chan types.Transaction)
		rc := make(chan types.Receipt)
		tw.wg.Add(1)

		go func(done <-chan struct{}) {
			defer tw.wg.Done()
			for {
				//if ok := tw.getTransactions(int64(i), tc); ok {
				if ok := tw.getTransactions2(int64(i), tc, rc); ok {
					i += int64(tw.workerCnt)
				}
				select {
				case <-done:
					tw.lastBlocks = append(tw.lastBlocks, int64(i))
					close(tc)
					close(rc)
					return
				default:
					time.Sleep(time.Millisecond)
				}
			}
		}(tw.done)

		return tc, rc
	}

	//merge := func(tcs ...<-chan types.Transaction) <-chan types.Transaction {
	merge := func(tcs []<-chan types.Transaction, rcs []<-chan types.Receipt) (<-chan types.Transaction, <-chan types.Receipt) {
		var wg sync.WaitGroup
		tout := make(chan types.Transaction, 100)
		rout := make(chan types.Receipt, 100)

		wg.Add(len(tcs))
		wg.Add(len(rcs))

		tsender := func(c <-chan types.Transaction) {
			for n := range c {
				tout <- n
			}

			wg.Done()
		}

		rsender := func(c <-chan types.Receipt) {
			for n := range c {
				rout <- n
			}

			wg.Done()
		}

		//		sender := func(tc <-chan types.Transaction, rc <-chan types.Receipt) {
		//		ENDSENDER:
		//			for {
		//				select {
		//				case t, ok := <-tc:
		//					if ok {
		//						tout <- t
		//					} else {
		//						tc = nil
		//					}
		//				case r, ok := <-rc:
		//					if ok {
		//						rout <- r
		//					} else {
		//						rc = nil
		//					}
		//				}
		//
		//				if tc == nil && rc == nil {
		//					break ENDSENDER
		//				}
		//			}
		//		}

		for _, tc := range tcs {
			go tsender(tc)
		}

		for _, rc := range rcs {
			go rsender(rc)
		}

		go func() {
			wg.Wait()
			close(tout)
			close(rout)
		}()

		return tout, rout
	}

	//wc := make([]<-chan types.Transaction, 0)
	tc := make([]<-chan types.Transaction, 0)
	rc := make([]<-chan types.Receipt, 0)

	for i := 0; i < tw.workerCnt; i++ {
		t, r := worker(tw.lastBlocks[i])
		tc = append(tc, t)
		rc = append(rc, r)
		//wc = append(wc, worker(tw.lastBlocks[i]))
	}

	tw.lastBlocks = make([]int64, 0)

	saveTxns := func(tc <-chan types.Transaction, rc <-chan types.Receipt) {
		tw.wg.Add(1)
		defer tw.wg.Done()

		db := db.GetRDBInstance()
		defer db.Close()

		//for txn := range tc {
		//	db.Insert(&txn)
		//}

	ENDSAVE:
		for {
			select {
			case txn, ok := <-tc:
				if ok {
					db.Insert(&txn)
				} else {
					tc = nil
				}
			case rpt, ok := <-rc:
				if ok {
					db.Insert(&rpt)
				} else {
					rc = nil
				}
			}

			if tc == nil && rc == nil {
				break ENDSAVE
			}
		}
	}

	//go saveTxns(merge(wc...))
	go saveTxns(merge(tc, rc))
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

func createRptTable() {
	db := db.GetRDBInstance()

	if err := db.CreateTable(&types.Receipt{}, &orm.CreateTableOptions{}); err == nil {
		fmt.Println("create receipt table successful")
	} else {
		return
	}

	db.Model((*types.Receipt)(nil)).Exec(`
		create index concurrently txnhash_idx on rpt_tbl("transactionHash")
	`)
}

func createTxnTable() {
	db := db.GetRDBInstance()

	if err := db.CreateTable(&types.Transaction{}, &orm.CreateTableOptions{}); err == nil {
		fmt.Println("create transaction table successful")
	} else {
		return
	}

	db.Model((*types.Transaction)(nil)).Exec(`
		create index concurrently hash_idx on txn_tbl(hash)
	`)
}

func NewTxnWorker(startBlock int64, url string, workerCnt int) *TxnWorker {
	createTxnTable()
	createRptTable()

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
