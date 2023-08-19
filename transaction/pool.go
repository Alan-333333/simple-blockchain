package transaction

type TxPool struct {
	Txs []*Transaction
}

var txPoolInstance *TxPool

// 创建交易池
func NewTxPool() *TxPool {

	if txPoolInstance != nil {
		return txPoolInstance
	}
	txPoolInstance = &TxPool{
		Txs: make([]*Transaction, 0),
	}

	return txPoolInstance
}

func GetTxPool() *TxPool {
	return txPoolInstance
}

// 添加新交易
func (pool *TxPool) AddTx(tx *Transaction) {
	pool.Txs = append(pool.Txs, tx)
}

// 从池中获取交易
func (pool *TxPool) GetTx() *Transaction {
	// 后入先出
	tx := pool.Txs[len(pool.Txs)-1]
	pool.Txs = pool.Txs[:len(pool.Txs)-1]
	return tx
}

// 获取交易池
func (pool *TxPool) GetTxs() []*Transaction {
	return pool.Txs
}

// 获取交易池大小
func (pool *TxPool) Size() int {
	return len(pool.Txs)
}

// 判断是否包含该交易
func (pool *TxPool) Has(tx *Transaction) bool {
	// 遍历查找
	for _, t := range pool.Txs {
		if t == tx {
			return true
		}
	}
	return false
}

// PopTransactions 从交易池中弹出指定数量的交易
func (pool *TxPool) PopTransactions(n int) []*Transaction {

	if n > len(pool.Txs) {
		// 交易池中交易不足
		return nil
	}

	txs := pool.Txs[:n]
	pool.Txs = pool.Txs[n:]

	return txs
}

// RemoveTransactions 从交易池中移除指定的交易
func (pool *TxPool) RemoveTransactions(txs []*Transaction) {

	result := []*Transaction{}
	for _, tx := range pool.Txs {
		removed := false
		for _, rmTx := range txs {
			if tx == rmTx {
				removed = true
				break
			}
		}
		if !removed {
			result = append(result, tx)
		}
	}

	pool.Txs = result
}
