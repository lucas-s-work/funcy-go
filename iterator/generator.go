package iterator

type NaturalGenerator struct {
	index int
}

func NewNaturalGenerator() Iterator[int] {
	return &NaturalGenerator{}
}

func (n *NaturalGenerator) Next() (int, error, bool) {
	n.index++
	return n.index - 1, nil, true
}

func (n *NaturalGenerator) Reset() error {
	n.index = 0
	return nil
}

type FibonnacciGenerator struct {
	a, b int
}

func NewFibonnacciGenerator() Iterator[int] {
	return &FibonnacciGenerator{
		a: 1,
		b: 0,
	}
}

func (f *FibonnacciGenerator) Next() (int, error, bool) {
	c := f.a
	f.a = f.a + f.b
	f.b = c

	return c, nil, true
}

func (f *FibonnacciGenerator) Reset() error {
	f.a = 1
	f.b = 0

	return nil
}

type MaskGenerator struct {
	stride int
	index  int
}

func NewMaskGenerator(stride int) Iterator[bool] {
	return &MaskGenerator{
		stride: stride,
		index:  stride,
	}
}

func (m *MaskGenerator) Next() (bool, error, bool) {
	pass := m.stride == m.index

	m.index++
	if pass {
		m.index = 0
	}

	return pass, nil, true
}

func (m *MaskGenerator) Reset() error {
	m.index = m.stride
	return nil
}
