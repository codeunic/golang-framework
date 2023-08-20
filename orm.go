package framework

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"reflect"
	"strings"
)

type Condition struct {
	Statement string
	Args      []interface{}
}

type Group struct {
	Conditions []Condition
	Operator   string
}
type Join struct {
	JoinType string
	Table    string
	On       string
}
type QueryBuilder struct {
	db         *sql.DB
	table      string
	fields     []string
	subQueries []string
	groups     []Group
	joins      []Join
	groupBy    []string
	order      string
	limit      int
	offset     int
	tx         *sql.Tx
}

func NewQueryBuilder(db *sql.DB, table string) *QueryBuilder {
	return &QueryBuilder{
		db:    db,
		table: table,
	}
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.fields = append(qb.fields, fields...)
	return qb
}

func (qb *QueryBuilder) Subquery(subquery string) *QueryBuilder {
	qb.subQueries = append(qb.subQueries, subquery)
	return qb
}

func (qb *QueryBuilder) Where(statement string, args ...interface{}) *QueryBuilder {
	if len(qb.groups) == 0 {
		qb.groups = append(qb.groups, Group{Operator: "AND"})
	}
	group := &qb.groups[len(qb.groups)-1]
	group.Conditions = append(group.Conditions, Condition{Statement: statement, Args: args})
	return qb
}

func (qb *QueryBuilder) Group(operator string) *QueryBuilder {
	qb.groups = append(qb.groups, Group{Operator: operator})
	return qb
}

func (qb *QueryBuilder) Order(order string) *QueryBuilder {
	qb.order = order
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder) Join(joinType, table, on string) *QueryBuilder {
	qb.joins = append(qb.joins, Join{
		JoinType: joinType,
		Table:    table,
		On:       on,
	})
	return qb
}

func (qb *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	qb.groupBy = fields
	return qb
}

func (qb *QueryBuilder) Begin() error {
	tx, err := qb.db.Begin()
	if err != nil {
		return err
	}
	qb.tx = tx
	return nil
}

func (qb *QueryBuilder) Commit() error {
	err := qb.tx.Commit()
	qb.tx = nil
	return err
}

func (qb *QueryBuilder) Rollback() error {
	err := qb.tx.Rollback()
	qb.tx = nil
	return err
}
func (qb *QueryBuilder) CreateIndex(table, field string) error {
	query := fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s)", table, field, table, field)
	_, err := qb.db.Exec(query)
	return err
}
func (qb *QueryBuilder) Build() (string, []interface{}) {
	var query strings.Builder
	var args []interface{}

	if len(qb.subQueries) > 0 {
		for _, subquery := range qb.subQueries {
			query.WriteString(fmt.Sprintf(" %s", subquery))
		}

		return query.String(), args
	}

	query.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(qb.fields, ", "), qb.table))

	for _, join := range qb.joins {
		query.WriteString(fmt.Sprintf(" %s JOIN %s ON %s", join.JoinType, join.Table, join.On))
	}

	for _, group := range qb.groups {
		var conditions []string
		for _, condition := range group.Conditions {
			conditions = append(conditions, condition.Statement)
			args = append(args, condition.Args...)
		}
		query.WriteString(fmt.Sprintf(" WHERE (%s)", strings.Join(conditions, fmt.Sprintf(" %s ", group.Operator))))
	}

	if qb.order != "" {
		query.WriteString(fmt.Sprintf(" ORDER BY %s", qb.order))
	}

	if qb.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	if qb.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return query.String(), args
}

func (qb *QueryBuilder) Query() (*sql.Rows, error) {
	query, args := qb.Build()
	return qb.db.Query(query, args...)
}

func (qb *QueryBuilder) ScanStructSlice(dest interface{}) error {
	sliceType := reflect.TypeOf(dest)
	isSlice := sliceType.Elem().Kind() == reflect.Slice

	rows, err := qb.Query()

	if err != nil {
		return err
	}

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
			fmt.Println(err)
			return err
		}

		// Crear un map para almacenar los valores mapeados a los nombres de columna
		row := make(map[string]interface{})

		for i, col := range columns {
			val := *values[i].(*interface{})
			switch v := val.(type) {
			case []uint8:
				var data interface{}
				fmt.Println()
				err := json.Unmarshal(v, &data)

				if err == nil {
					row[col] = data
					continue
				}

				var strArray pq.StringArray
				err = strArray.Scan(v)

				if err == nil {
					row[col] = strArray
				} else {
					row[col] = val
				}
			default:
				row[col] = val
			}
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
