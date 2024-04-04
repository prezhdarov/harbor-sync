package harbor

type Label struct {
	ID           int    `json:"id"`
	Color        string `json:"color"`
	CreationTime string `json:"creation_time"`
	Name         string `json:"name"`
	Scope        string `json:"scope"`
	UpdateTime   string `json:"update_time"`
}

type Harbor struct {
	Host *string
	User *string
	Pass *string
}
