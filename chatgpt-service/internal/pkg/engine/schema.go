package engine

type Schema struct {
	ChainName   string
	SchemaName  string
	TableName   string
	TableSchema string
}

// EthereumCore : The database schema for EthereumCore blockchain
const (
	Ethereum        = "ethereum"
	CoreSchema      = "core"
	FactTransaction = "fact.transaction"
)

func (sh *Schema) String() string {
	// ethereum.core.fact.transaction
	tableName := sh.ChainName + "." + sh.SchemaName + "." + sh.TableName
	message := "Given the table below.\n\ntable name: " + tableName
	message += "\n\ntable schema: " + sh.TableSchema
	return message
}

func CreateEthereumCoreTransactionSchema() *Schema {
	ChainName := Ethereum
	SchemaName := CoreSchema
	TableName := FactTransaction
	// TODO: Schema to have different struct
	TableSchema := "BLOCK_HASH text, BLOCK_NUMBER number, BLOCK_TIMESTAMP timestamp_ntz, CUMULATIVE_GAS_USED float, ETH_VALUE float, FROM_ADDRESS text, GAS_LIMIT number, GAS_PRICE float, GAS_USED float"
	return &Schema{ChainName: ChainName, SchemaName: SchemaName, TableName: TableName, TableSchema: TableSchema}
}
