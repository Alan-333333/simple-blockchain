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
