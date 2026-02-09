package repository

import _ "embed"

//go:embed sql/events/create.sql
var qEventCreate string

//go:embed sql/events/get_all.sql
var qEventGetAll string

//go:embed sql/events/get_by_id.sql
var qEventGetByID string

//go:embed sql/events/update.sql
var qEventUpdate string

//go:embed sql/events/delete.sql
var qEventDelete string
