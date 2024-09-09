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

		collection, err := dao.FindCollectionByNameOrId("ktiu1ux24jyuxh7")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("1wpzyevj")

		// remove
		collection.Schema.RemoveField("oq1vws3u")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ktiu1ux24jyuxh7")
		if err != nil {
			return err
		}

		// add
		del_latitude := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "1wpzyevj",
			"name": "latitude",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), del_latitude); err != nil {
			return err
		}
		collection.Schema.AddField(del_latitude)

		// add
		del_longitude := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "oq1vws3u",
			"name": "longitude",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), del_longitude); err != nil {
			return err
		}
		collection.Schema.AddField(del_longitude)

		return dao.SaveCollection(collection)
	})
}
