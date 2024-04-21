package main

type ColumnType string

const (
	ColumnTypeInt    ColumnType = "int"
	ColumnTypeBool   ColumnType = "bool"
	ColumnTypeFloat  ColumnType = "float"
	ColumnTypeString ColumnType = "string"
)

type Table struct {
	ID         int      `json:"id"`
	DatabaseID int      `json:"database_id"`
	Name       string   `json:"name"`
	Columns    []Column `json:"columns"`
}

type Column struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Type       ColumnType `json:"type"`
	Increments bool       `json:"increments"`
}
