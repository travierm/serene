package main

type Database struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

var databases = []Database{
	{ID: 1, Name: "serene_main", Tables: []Table{
		{ID: 1, DatabaseID: 1, Name: "products", Columns: []Column{
			{ID: 1, Name: "id", Type: ColumnTypeInt, Increments: true},
			{ID: 2, Name: "name", Type: ColumnTypeString, Increments: false},
			{ID: 3, Name: "price", Type: ColumnTypeFloat, Increments: false},
		}},
	},
	}}
