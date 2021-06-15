package api

import kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"

type API struct {
	// We could internally extend this struct with specific Kamrad needs
	// or cherry pick the Kuadrant ones.
	Name string
	Spec kctlrv1beta1.APISpec
}
