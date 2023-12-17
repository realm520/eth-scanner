package transaction

import (
	"fmt"
	"strconv"
	"strings"
	"log"
	"os"
	"sync"
	"time"
)

type TransactionReporter struct {
	outputFile           *os.File
	filteredTransactions chan *TransactionResult
	waitGroup            *sync.WaitGroup
	done                 bool
}

func NewTransactionReporter(outputFile *os.File, transactions chan *TransactionResult, wg *sync.WaitGroup) *TransactionReporter {
	return &TransactionReporter{
		outputFile:           outputFile,
		filteredTransactions: transactions,
		waitGroup:            wg,
	}
}

func hex2dec(hexString string) int64 {
    hexString = strings.TrimPrefix(hexString, "0x")
    decimalValue, err := strconv.ParseInt(hexString, 16, 64)
    if err != nil {
		fmt.Println("转换失败:", err)
		return -1
	} else {
        return decimalValue
    }
}

func (reporter *TransactionReporter) Start() {
	defer reporter.waitGroup.Done()

	for !reporter.done {
		select {
		case transaction := <-reporter.filteredTransactions:
			transactionRow := fmt.Sprintf("%s,%s,%s,%d,%d,%d,%s\n", transaction.Hash, transaction.From, transaction.To, hex2dec(transaction.BlockNumber), hex2dec(transaction.TransactionIndex), 0, transaction.Input)
			if _, err := reporter.outputFile.WriteString(transactionRow); err != nil {
				log.Println("Error writing transaction to file:", err.Error())
				reporter.done = true
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (reporter *TransactionReporter) Stop() {
	reporter.done = true
}
