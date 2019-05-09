package api 

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/septemhill/ethacctdb/db"
	"github.com/septemhill/ethacctdb/types"
)

const maxTxns = 100

func GetAccountTotalTxnsCount(ctx *gin.Context) {
	d := db.GetRDBInstance()
	var cnt int

	d.Query(&cnt, fmt.Sprintf("select count(*) from txn_tbl where txn_from = '%s' or txn_to = '%s'", ctx.Param("addr"), ctx.Param("addr")))

	ctx.JSON(http.StatusOK, cnt)
}

func GetAccountTxns(ctx *gin.Context) {
	d := db.GetRDBInstance()
	txns := make([]types.Transaction, 0)
	
	limit := ctx.DefaultQuery("limit", "20")
	offset := ctx.DefaultQuery("offset", "0")
	asc := ctx.DefaultQuery("asc", "false")

	v, _ := strconv.ParseBool(asc)

	if !v {
		d.Query(&txns, fmt.Sprintf("select txn_tbl.*, rpt_tbl.status from txn_tbl inner join rpt_tbl on txn_tbl.hash = rpt_tbl.\"transactionHash\" where txn_from = '%s' or txn_to = '%s' order by ts desc limit %s offset %s", ctx.Param("addr"), ctx.Param("addr"), limit, offset))
	} else {
		d.Query(&txns, fmt.Sprintf("select txn_tbl.*, rpt_tbl.status from txn_tbl inner join rpt_tbl on txn_tbl.hash = rpt_tbl.\"transactionHash\" where txn_from = '%s' or txn_to = '%s' order by ts limit %s offset %s", ctx.Param("addr"), ctx.Param("addr"), limit, offset))
	}

	ctx.JSON(http.StatusOK, txns)
}

func GetHashInfo(ctx *gin.Context) {
	d := db.GetRDBInstance()

	filter := ctx.DefaultQuery("filter", "all")

	if filter == "all" || filter == "txn" {
		var txn *types.Transaction
		_, err := d.Query(&txn, fmt.Sprintf("select * from txn_tbl where hash = '%s'", ctx.Param("hash")))
	
		if err == nil && txn != nil {
			ctx.JSON(http.StatusOK, txn)
		}
		return
	}

	if filter == "all" || filter == "contract" {
		var rpt *types.Receipt
		_, err := d.Query(&rpt, fmt.Sprintf("select * from rpt_tbl where contractAddress = '%s'", ctx.Param("hash")))

		if err == nil && rpt != nil {
			ctx.JSON(http.StatusOK, rpt)
		}
		return
	}
}