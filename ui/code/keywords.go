package code

// TokenType represents the type of SQL token
type TokenType int

const (
	// Token types
	Statement TokenType = iota
	Object
	Join
	Operator
	Value
	Modifier
	Clause
	Sort
	Function
	Expression
	SetOperator
	CTE
	Values
	Constraint
	Attribute
	Transaction
	Permission
	Action
	Control
	Result
	Meta
	DataType
	Analytics
)

// Mapping of token types to their string representation (for debugging/display)
var TokenTypeNames = map[TokenType]string{
	Statement:   "Statement",
	Object:      "Object",
	Join:        "Join",
	Operator:    "Operator",
	Value:       "Value",
	Modifier:    "Modifier",
	Clause:      "Clause",
	Sort:        "Sort",
	Function:    "Function",
	Expression:  "Expression",
	SetOperator: "SetOperator",
	CTE:         "CTE",
	Values:      "Values",
	Constraint:  "Constraint",
	Attribute:   "Attribute",
	Transaction: "Transaction",
	Permission:  "Permission",
	Action:      "Action",
	Control:     "Control",
	Result:      "Result",
	Meta:        "Meta",
	DataType:    "DataType",
	Analytics:   "Analytics",
}

// Keyword represents an SQL keyword with its type
type Keyword struct {
	Token string
	Type  TokenType
}

// SQLKeywords is a list of SQL keywords that should be highlighted
var SQLKeywords = []Keyword{
	// Statement keywords
	{Token: "SELECT", Type: Statement},
	{Token: "FROM", Type: Statement},
	{Token: "WHERE", Type: Statement},
	{Token: "INSERT", Type: Statement},
	{Token: "INTO", Type: Statement},
	{Token: "UPDATE", Type: Statement},
	{Token: "DELETE", Type: Statement},
	{Token: "CREATE", Type: Statement},
	{Token: "DROP", Type: Statement},
	{Token: "ALTER", Type: Statement},
	{Token: "SET", Type: Statement},

	// Object keywords
	{Token: "TABLE", Type: Object},
	{Token: "VIEW", Type: Object},
	{Token: "INDEX", Type: Object},
	{Token: "DATABASE", Type: Object},
	{Token: "SCHEMA", Type: Object},
	{Token: "TRIGGER", Type: Object},
	{Token: "PROCEDURE", Type: Object},
	{Token: "FUNCTION", Type: Object},

	// Join keywords
	{Token: "JOIN", Type: Join},
	{Token: "INNER", Type: Join},
	{Token: "LEFT", Type: Join},
	{Token: "RIGHT", Type: Join},
	{Token: "FULL", Type: Join},
	{Token: "OUTER", Type: Join},
	{Token: "CROSS", Type: Join},
	{Token: "LATERAL", Type: Join},

	// Operator keywords
	{Token: "ON", Type: Operator},
	{Token: "AS", Type: Operator},
	{Token: "AND", Type: Operator},
	{Token: "OR", Type: Operator},
	{Token: "NOT", Type: Operator},
	{Token: "IN", Type: Operator},
	{Token: "BETWEEN", Type: Operator},
	{Token: "LIKE", Type: Operator},
	{Token: "IS", Type: Operator},

	// Value keywords
	{Token: "NULL", Type: Value},
	{Token: "TRUE", Type: Value},
	{Token: "FALSE", Type: Value},
	{Token: "DEFAULT", Type: Value},

	// Modifier keywords
	{Token: "DISTINCT", Type: Modifier},
	{Token: "ALL", Type: Modifier},
	{Token: "ANY", Type: Modifier},

	// Clause keywords
	{Token: "GROUP", Type: Clause},
	{Token: "HAVING", Type: Clause},
	{Token: "ORDER", Type: Clause},
	{Token: "BY", Type: Clause},
	{Token: "LIMIT", Type: Clause},
	{Token: "OFFSET", Type: Clause},

	// Sort keywords
	{Token: "ASC", Type: Sort},
	{Token: "DESC", Type: Sort},

	// Function keywords
	{Token: "COUNT", Type: Function},
	{Token: "SUM", Type: Function},
	{Token: "AVG", Type: Function},
	{Token: "MIN", Type: Function},
	{Token: "MAX", Type: Function},
	{Token: "COALESCE", Type: Function},
	{Token: "CONCAT", Type: Function},
	{Token: "SUBSTRING", Type: Function},
	{Token: "UPPER", Type: Function},
	{Token: "LOWER", Type: Function},
	{Token: "TRIM", Type: Function},
	{Token: "CURRENT_DATE", Type: Function},
	{Token: "CURRENT_TIME", Type: Function},
	{Token: "CURRENT_TIMESTAMP", Type: Function},

	// Expression keywords
	{Token: "CASE", Type: Expression},
	{Token: "WHEN", Type: Expression},
	{Token: "THEN", Type: Expression},
	{Token: "ELSE", Type: Expression},
	{Token: "END", Type: Expression},

	// Set operation keywords
	{Token: "UNION", Type: SetOperator},
	{Token: "INTERSECT", Type: SetOperator},
	{Token: "EXCEPT", Type: SetOperator},

	// CTE keywords
	{Token: "WITH", Type: CTE},

	// Values keywords
	{Token: "VALUES", Type: Values},

	// Constraint keywords
	{Token: "PRIMARY", Type: Constraint},
	{Token: "FOREIGN", Type: Constraint},
	{Token: "KEY", Type: Constraint},
	{Token: "CONSTRAINT", Type: Constraint},
	{Token: "REFERENCES", Type: Constraint},
	{Token: "UNIQUE", Type: Constraint},
	{Token: "CHECK", Type: Constraint},

	// Attribute keywords
	{Token: "AUTO_INCREMENT", Type: Attribute},

	// Transaction keywords
	{Token: "BEGIN", Type: Transaction},
	{Token: "COMMIT", Type: Transaction},
	{Token: "ROLLBACK", Type: Transaction},
	{Token: "TRANSACTION", Type: Transaction},

	// Permission keywords
	{Token: "GRANT", Type: Permission},
	{Token: "REVOKE", Type: Permission},
	{Token: "TO", Type: Permission},

	// Action keywords
	{Token: "CASCADE", Type: Action},
	{Token: "RESTRICT", Type: Action},

	// Control keywords
	{Token: "IF", Type: Control},
	{Token: "EXISTS", Type: Control},

	// Result keywords
	{Token: "RETURNING", Type: Result},

	// Meta keywords
	{Token: "EXPLAIN", Type: Meta},
	{Token: "ANALYZE", Type: Meta},

	// Data type keywords
	{Token: "DATE", Type: DataType},
	{Token: "TIME", Type: DataType},
	{Token: "TIMESTAMP", Type: DataType},
	{Token: "INTERVAL", Type: DataType},
	{Token: "INT", Type: DataType},
	{Token: "SERIAL", Type: DataType},
	{Token: "INTEGER", Type: DataType},
	{Token: "SMALLINT", Type: DataType},
	{Token: "BIGINT", Type: DataType},
	{Token: "DECIMAL", Type: DataType},
	{Token: "NUMERIC", Type: DataType},
	{Token: "FLOAT", Type: DataType},
	{Token: "REAL", Type: DataType},
	{Token: "DOUBLE", Type: DataType},
	{Token: "CHAR", Type: DataType},
	{Token: "VARCHAR", Type: DataType},
	{Token: "TEXT", Type: DataType},
	{Token: "BOOLEAN", Type: DataType},
	{Token: "BINARY", Type: DataType},
	{Token: "BLOB", Type: DataType},
	{Token: "JSON", Type: DataType},

	// Analytics keywords
	{Token: "WINDOW", Type: Analytics},
	{Token: "PARTITION", Type: Analytics},
	{Token: "OVER", Type: Analytics},
	{Token: "RANK", Type: Analytics},
	{Token: "ROW_NUMBER", Type: Analytics},
	{Token: "DENSE_RANK", Type: Analytics},
	{Token: "NTILE", Type: Analytics},
	{Token: "LEAD", Type: Analytics},
	{Token: "LAG", Type: Analytics},
}

// GetKeywordMap returns a map of SQL keywords for faster lookup
func GetKeywordMap() map[string]TokenType {
	keywordMap := make(map[string]TokenType, len(SQLKeywords))
	for _, keyword := range SQLKeywords {
		keywordMap[keyword.Token] = keyword.Type
	}
	return keywordMap
}
