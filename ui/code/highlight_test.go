package code

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHighlight(t *testing.T) {
	tests := []string{
		"SELECT id, name FROM users WHERE age > 18;",
		"INSERT INTO products(name, price) VALUES('Laptop', 999.99);",
		"CREATE TABLE employees(id INT PRIMARY KEY, name VARCHAR(100), salary DECIMAL(10,2));",
		"UPDATE orders SET status='shipped' WHERE order_date < '2023-01-01' AND total > 100;",
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fmt.Println(Highlight(tt))
		})
	}
}
