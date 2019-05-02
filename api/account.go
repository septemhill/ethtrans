package api 

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/septemhill/ethacctdb/db"
	"github.com/septemhill/ethacctdb/types"
)

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

	d.Query(&txns, fmt.Sprintf("select * from txn_tbl where txn_from = '%s' or txn_to = '%s' order by ts limit %s offset %s", ctx.Param("addr"), ctx.Param("addr"), limit, offset))

	ctx.JSON(http.StatusOK, txns)
}

