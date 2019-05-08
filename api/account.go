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
		d.Query(&txns, fmt.Sprintf("select * from txn_tbl inner join rpt_tbl on txn_tbl.hash = rpt_tbl.\"transactionHash\" where txn_from = '%s' or txn_to = '%s' order by ts desc limit %s offset %s", ctx.Param("addr"), ctx.Param("addr"), limit, offset))
	} else {
		d.Query(&txns, fmt.Sprintf("select * from txn_tbl inner join rpt_tbl on txn_tbl.hash = rpt_tbl.\"transactionHash\" where txn_from = '%s' or txn_to = '%s' order by ts limit %s offset %s", ctx.Param("addr"), ctx.Param("addr"), limit, offset))
	}

	ctx.JSON(http.StatusOK, txns)
}

