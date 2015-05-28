package tree

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
)

type BaseNode struct {
	id       string
	children map[string]Node
	tags     map[string]interface{}
}

func NewBaseNode(kv map[string]interface{}) (bn *BaseNode, err error) {
	for _, v := range kv {
		switch v.(type) {
		case uint64, float64, int64, string:
		default:
			err = fmt.Errorf("Value %v must be uint64, int64, float64 or string", v)
			return
		}
	}
	bn = &BaseNode{
		id:       uuid.New(),
		tags:     kv,
		children: make(map[string]Node, 4),
	}
	return
}

func InitBaseNode(bn *BaseNode, kv map[string]interface{}) (err error) {
	for _, v := range kv {
		switch v.(type) {
		case uint64, float64, int64, string:
		default:
			err = fmt.Errorf("Value %v must be uint64, int64, float64 or string", v)
			return
		}
	}
	bn.id = uuid.New()
	if kv == nil {
		bn.tags = make(map[string]interface{})
	} else {
		bn.tags = kv
	}
	bn.children = make(map[string]Node, 4)
	return
}

func (bn *BaseNode) Id() string {
	return bn.id
}

func (bn *BaseNode) Children() map[string]Node {
	return bn.children
}

func (bn *BaseNode) AddChild(n Node) bool {
	var found bool
	if _, found = bn.children[n.Id()]; !found {
		bn.children[n.Id()] = n
	}
	return found
}

func (bn *BaseNode) Input(args ...interface{}) error {
	return fmt.Errorf("BaseNode has no Input")
}

func (bn *BaseNode) Output() (interface{}, error) {
	fmt.Printf("Node kv: %v\n", bn.tags)
	return nil, fmt.Errorf("BaseNode has no Output")
}

func (bn *BaseNode) Get(key string) (val interface{}, found bool) {
	val, found = bn.tags[key]
	return
}

func (bn *BaseNode) Set(key string, val interface{}) {
	bn.tags[key] = val
}
