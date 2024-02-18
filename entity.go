package packme

type ItemInput struct {
	Id  string `json:"id"`
	Qty int    `json:"qty"`
	Dim [3]int `json:"dim"`
}

type PackmeInput struct {
	Containers []ItemInput `json:"containers"`
	Items      []ItemInput `json:"items"`
}

type Rotation string // "LWH" | "WLH" | "WHL" | "HLW" | "HWL" | "LHW"

type Vector3 struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type ItemResult struct {
	Id  string   `json:"id"`
	Dim Vector3  `json:"dim"`
	Pos Vector3  `json:"pos"`
	Rot Rotation `json:"rot"`
}

type ContainerOutput struct {
	Id    string       `json:"id"`
	Dim   Vector3      `json:"dim"`
	Items []ItemResult `json:"items"`
}

type PackmeOutput struct {
	Containers    []ContainerOutput `json:"containers"`
	UnpackedItems []ItemResult      `json:"unpacked_items"`
}
