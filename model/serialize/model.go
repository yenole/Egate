//+build serialize

package serialize

type Model struct {
	Id uint64 `gorm:"primary_key"` // 主键
}

func (p *Model) ID() uint64 {
	return p.Id
}

func (p *Model) SetId(v uint64) {
	p.Id = v
}
