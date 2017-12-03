package pager

import (
	"net/http"
	"net/url"
	"strconv"

	utils2 "github.com/go-qbit/template/utils"
)

type Pager struct {
	curPage       uint64
	perPage       uint64
	r             *http.Request
	totalElements uint64
}

func NewPager(r *http.Request, perPage uint64) *Pager {
	curPage, err := strconv.ParseUint(r.Form.Get("page"), 10, 64)
	if err != nil {
		curPage = 1
	}

	if curPage < 1 {
		curPage = 1
	}

	return &Pager{
		curPage: curPage,
		perPage: perPage,
		r:       r,
	}
}

func (p *Pager) CurPage() uint64 {
	return p.curPage
}

func (p *Pager) PerPage() uint64 {
	return p.perPage
}

func (p *Pager) SetTotalElements(v uint64) {
	p.totalElements = v
}

func (p *Pager) TotalElements() uint64 {
	return p.totalElements
}

func (p *Pager) TotalPages() uint64 {
	if p.totalElements == 0 {
		return 0
	}

	res := p.totalElements / p.perPage
	if p.totalElements%p.perPage != 0 {
		res++
	}

	return res
}

func (p *Pager) GetLink(page uint64, overrideValues ...interface{}) string {
	newParams := url.Values{}
	for k, v := range p.r.Form {
		newParams[k] = v
	}

	i := 0
	for i < len(overrideValues) {
		if i+1 < len(overrideValues) {
			newParams.Set(utils2.ToString(overrideValues[i]), utils2.ToString(overrideValues[i+1]))
		} else {
			newParams.Set(utils2.ToString(overrideValues[i]), "")
		}
		i += 2
	}

	newParams.Set("page", strconv.FormatUint(page, 10))

	return "?" + newParams.Encode()
}

func (p *Pager) GetOffset() uint64 {
	return (p.curPage - 1) * p.perPage
}
