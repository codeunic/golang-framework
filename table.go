package framework

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"reflect"
	"strings"
)

// ... other definitions ...

type Table struct {
	db        *sql.DB
	tableName string
}

// NewTable
// users := orm.NewTable(db, "users")
//
//	// Create
//	data := map[string]interface{}{
//		"name": "John Doe",
//		"age":  30,
//	}
//	result, err := users.Create(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Read
//	rows, err := users.Read("age > ?", 18)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Process rows here...
//
//	// Update
//	data = map[string]interface{}{
//		"age": 31,
//	}
//	result, err = users.Update```go
//	// Update
//	data = map[string]interface{}{
//		"age": 31,
//	}
//	result, err = users.Update("name = ?", data, "John Doe")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Delete
//	result, err = users.Delete("age < ?", 18)
//	if err != nil {
//		log.Fatal(err)
//	}
func NewTable(db *sql.DB, tableName string) *Table {
	return &Table{
		db:        db,
		tableName: tableName,
	}
}

func (t *Table) Create(data map[string]interface{}) (sql.Result, error) {
	fields := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for field, value := range data {
		fields = append(fields, field)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		t.tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)

	return t.db.Exec(query, values...)
}

func (t *Table) Read(columns string, where string, args ...interface{}) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", columns, t.tableName, where)
	return t.db.Query(query, args...)
}

func (t *Table) ScanStructSlice(rows *sql.Rows, dest interface{}) error {
	sliceType := reflect.TypeOf(dest)
	isSlice := sliceType.Elem().Kind() == reflect.Slice

	var results []interface{}

	// Obtener los nombres de los campos del resultado
	columns, err := rows.Columns()

	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))

	for i := range values {
		values[i] = new(interface{})
	}

	// Escanear los valores en el slice de interfaces
	for rows.Next() {
		err := rows.Scan(values...)

		if err != nil {
			return err
		}

		// Crear un map para almacenar los valores mapeados a los nombres de columna
		row := make(map[string]interface{})

		for i, col := range columns {
			row[col] = *values[i].(*interface{})
		}

		// Agregar el map al slice de resultados
		results = append(results, row)
	}

	var marshal []byte

	if isSlice {
		marshal, err = json.Marshal(results)
	} else {
		if len(results) == 0 {
			return errors.New("no se ha encontrado coincidencias")
		}
		marshal, err = json.Marshal(results[0])
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(marshal, &dest)
}

func (t *Table) Update(where string, data map[string]interface{}, args ...interface{}) (sql.Result, error) {
	setParts := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data)+len(args))
	for field, value := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", field))
		values = append(values, value)
	}
	values = append(values, args...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		t.tableName,
		strings.Join(setParts, ", "),
		where,
	)

	return t.db.Exec(query, values...)
}

func (t *Table) Delete(where string, args ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", t.tableName, where)
	return t.db.Exec(query, args...)
}
