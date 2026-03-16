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

//go:embed sql/users/create.sql
var qUserCreate string

//go:embed sql/users/get_by_email.sql
var qUserGetByEmail string

//go:embed sql/users/get_by_id.sql
var qUserGetByID string
