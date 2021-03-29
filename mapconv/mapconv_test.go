package mapconv

type V struct {
	B bool              `json:"bool"`
	I int               `json:"int"`
	S string            `json:"string"`
	A []int             `json:"array_int"`
	M map[string]string `json:"map_string"`
	P P                 `json:"p"`
	Q Q                 `json:"q"`
	R int               `json:"-"`
}

type P struct {
	B *bool              `json:"bool"`
	I *int               `json:"int"`
	S *string            `json:"string"`
	A []*int             `json:"array_int"`
	M map[string]*string `json:"map_string"`
}

type Q struct {
	B bool
	I int
	S string
	A []int
	M map[string]string
}
