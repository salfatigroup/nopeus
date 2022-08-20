package config

import "encoding/gob"

func init() {
	// register interfaces with gob
	gob.Register(&IngressPath{})
	gob.Register(&Ingress{})
	gob.Register([]*IngressPath{})
	gob.Register([]*Ingress{})
	gob.Register(&NopeusDefaultMicroservice{})
	gob.Register(map[string]string{})
}
