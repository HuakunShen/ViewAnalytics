package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("g9evvmsbp5k8552")
		if err != nil {
			return err
		}

		// add
		new_type := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "kzdpbzqm",
			"name": "type",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"proxy",
					"redirect"
				]
			}
		}`), new_type); err != nil {
			return err
		}
		collection.Schema.AddField(new_type)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("g9evvmsbp5k8552")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("kzdpbzqm")

		return dao.SaveCollection(collection)
	})
}
