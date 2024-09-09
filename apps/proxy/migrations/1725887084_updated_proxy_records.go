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
		new_referer := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ayyjxpln",
			"name": "referer",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_referer); err != nil {
			return err
		}
		collection.Schema.AddField(new_referer)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("g9evvmsbp5k8552")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ayyjxpln")

		return dao.SaveCollection(collection)
	})
}
