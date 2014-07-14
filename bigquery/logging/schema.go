package logging

import (
	"fmt"
	"strings"
)

// map[テーブル名]map[列の名前]列の種類
var schemata = make(map[string]map[string]string)

// schemaを登録します。
func AddSchema(table, schema string) {
	if _, ok := schemata[table]; ok {
		panic(fmt.Errorf("table[%s] that already exists.", table))
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

		// column名、type名のいずれかが欠けているか。
		if err := confirmFormat(nameAndType); err != nil {
			panic(fmt.Errorf("%s\nYour input schema: %s", err.Error(), schema))
		}
		// 許可されたtypeか確認します。
		if err := confirmType(nameAndType); err != nil {
			panic(fmt.Errorf("%s", err.Error()))
		}

		schemaMap[nameAndType[0]] = strings.ToUpper(nameAndType[1])
	}

	schemata[table] = make(map[string]string)
	schemata[table] = schemaMap
}

func confirmFormat(nameAndType []string) (err error) {
	if len(nameAndType) != 2 {
		return fmt.Errorf("Format of the schema is invalid. ex. \"column1_name:data_type,column2_name:data_type,...\"")
	}
	return nil
}

func confirmType(nameAndType []string) (err error) {
	// 許可されたtypeか否か確認します。
	switch strings.ToUpper(nameAndType[1]) {
	case "STRING":
	case "INTEGER":
	case "FLOAT":
	case "BOOLEAN":
	case "TIMESTAMP":
	case "RECORD":
	default:
		return fmt.Errorf("Invalid Type: %s\nValid type: STRING, INTEGER, FLOAT, BOOLEAN, TIMESTAMP, RECORD", nameAndType[1])
	}
	return nil
}

/*
func (s *SchemaService) GetSchema(key string) string {
	return ""
}
*/
