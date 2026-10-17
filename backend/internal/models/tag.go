package models

// Tag is a coloured label that can be attached to tasks and deadlines.
type Tag struct {
	Meta  `bson:",inline"`
	Title string `bson:"title" json:"title"`
	// Color is a hex colour such as "#8CFBDE".
	Color string `bson:"color" json:"color"`
}
