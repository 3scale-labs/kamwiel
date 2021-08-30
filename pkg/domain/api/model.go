package api

type API struct {
	Name string `json:"name,omitempty"`
	// Inline OAS
	// +optional
	Spec string `json:"spec,omitempty"`
}

type APIs []API
