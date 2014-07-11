package logging

import (
	"fmt"
	"strings"
)

// map[テーブル名]map[列の名前]列の種類
var schemata = make(map[string]map[string]string)

type SchemaService struct{}

func (s *SchemaService) Init(definition map[string]string) {
	for table, schema := range definition {
		err := s.AddSchema(table, schema)
		if err != nil {
			panic(fmt.Errorf("%s", err.Error()))
		}
	}
}

// schemaを登録します。
func (s *SchemaService) AddSchema(table, schema string) (err error) {
	if _, ok := schemata[table]; ok {
		return fmt.Errorf("%s table that already exists.", table)
	}

	// コンマでschemaを分割します。["kind:string", "date:timestamp", "count:integer"]
	col_schemata := strings.FieldsFunc(schema, func(r rune) bool {
		return strings.ContainsRune(",", r)
	})

	schemaMap := make(map[string]string)
	for _, col_schema := range col_schemata {
		// さらにコロンで分割します。["kind", "string"]
		nameAndType := strings.FieldsFunc(col_schema, func(r rune) bool {
			return strings.ContainsRune(":", r)
		})

		// ここでcolumn名、type名が正しいか確認する。
		err := confirmNameAndType(nameAndType)
		if err != nil {
			return err
		}

		schemaMap[nameAndType[0]] = strings.ToUpper(nameAndType[1])
	}

	schemata[table] = make(map[string]string)
	schemata[table] = schemaMap
	return nil
}

func confirmNameAndType(nameAndType []string) (err error) {
	if len(nameAndType) != 2 {
		return fmt.Errorf("Format of the schema is invalid. ex. \"column1_name:data_type,column2_name:data_type,...\"")
	}

	name := nameAndType[0]
	typeName := nameAndType[1]

	// 列と種類が空か否か確認します。
	if len(name) == 0 {
		return fmt.Errorf("Column is empty.")
	}
	if len(typeName) == 0 {
		return fmt.Errorf("Type is empty.")
	}

	// 許可されたtypeか否か確認します。
	switch strings.ToUpper(typeName) {
	case "STRING":
	case "INTEGER":
	case "FLOAT":
	case "BOOLEAN":
	case "TIMESTAMP":
	case "RECORD":
	default:
		return fmt.Errorf("Invalid Type: %s\nValid type: STRING, INTEGER, FLOAT, BOOLEAN, TIMESTAMP, RECORD", typeName)
	}

	return nil
}

func (s *SchemaService) GetSchema(key string) string {
	return ""
}
