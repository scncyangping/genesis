package padding

import (
	"genesis/pkg/core/dandelion"
)

type Padding struct {
	paddings []dandelion.PaddingOptions
	// 用于填充模版的参数数据
	m map[string]any
}

func NewPadding() *Padding {
	return &Padding{m: map[string]any{}}
}

func (pa *Padding) AddPadding(p dandelion.PaddingOptions) {
	if pa.paddings == nil {
		pa.paddings = make([]dandelion.PaddingOptions, 0)
	}
	pa.paddings = append(pa.paddings, p)
}

func (pa *Padding) Padding() map[string]any {
	for _, v := range pa.paddings {
		for key, value := range v() {
			pa.m[key] = value
		}
	}
	return pa.m
}
