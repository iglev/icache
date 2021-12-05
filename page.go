package icache

import (
	"context"
	"fmt"
	"reflect"
)

// PageMeta page meta
type PageMeta struct {
	PassBack string
	Total    int
	HasMore  bool
}

// PageOpMeta page op meta
type PageOpMeta struct {
	Limit  int
	Offset int
}

// PageOp page op
type PageOp interface {
	GenNextPassBack(strPassBack string) string
	PasePassBack(strPassBack string) PageOpMeta
	Get(ctx context.Context, key string, pDestSink interface{}) error
}

// Page page
type Page struct {
	op PageOp
}

// NewPage new page
func NewPage(op PageOp) (*Page, error) {
	if op == nil {
		return nil, fmt.Errorf("invalid param")
	}
	return &Page{op: op}, nil
}

// Get get page
func (p *Page) Get(ctx context.Context, key string, strPassBack string, pDestSink interface{}) (PageMeta, error) {
	meta := PageMeta{}
	// 检查pDestSink是否为数组指针
	sinkType := reflect.TypeOf(pDestSink)
	if sinkType.Kind() != reflect.Ptr {
		return meta, fmt.Errorf("pDestSink not ptr type")
	}
	sinkVal := reflect.ValueOf(pDestSink)
	pSinkValElem := sinkVal.Elem()
	if pSinkValElem.Kind() != reflect.Slice {
		return meta, fmt.Errorf("pDestSink not point slice")
	}

	err := p.op.Get(ctx, key, pDestSink)
	if err != nil {
		return meta, err
	}

	meta.Total = pSinkValElem.Len()
	meta.PassBack = p.op.GenNextPassBack(strPassBack)
	opMeta := p.op.PasePassBack(strPassBack)
	if opMeta.Offset < 0 || opMeta.Limit < 0 {
		return meta, nil
	}
	end := opMeta.Offset + opMeta.Limit
	start := opMeta.Offset
	if start >= meta.Total {
		start = meta.Total
	}
	var tmpVal reflect.Value
	if meta.Total <= 0 {
		// 没有数据
		return meta, nil
	} else if end < meta.Total {
		// 还有更多
		meta.HasMore = true
		tmpVal = pSinkValElem.Slice(start, end)
	} else {
		// 没有更多
		tmpVal = pSinkValElem.Slice(start, meta.Total)
	}
	pSinkValElem.Set(tmpVal)
	return meta, nil
}
