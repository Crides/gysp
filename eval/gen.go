package eval

type Generator struct {
    sources     [][]*Object
    ptrs        []int
    sizes       []int
    size        int
}

func NewGenerator() *Generator {
    return new(Generator)
}

func (g * Generator) AddSource(a []*Object) {
    g.sources = append(g.sources, a)
    g.ptrs = append(g.ptrs, 0)
    g.sizes = append(g.sizes, len(a))
    g.size ++
}

func (g * Generator) Generate(f func([]*Object) *Object) {
    for g.ptrs[0] < g.sizes[0] {
        vals := make([]*Object, g.size)
        for i := 0; i < g.size; i ++ {
            vals[i] = g.sources[i][g.ptrs[i]]
        }
        f(vals)

        // Increment pointers
        g.ptrs[g.size - 1] ++       // Last pointer
        for i := g.size - 1; i > 0; i -- {
            if g.ptrs[i] >= g.sizes[i] {
                g.ptrs[i] = 0
                g.ptrs[i - 1] ++
            }
        }
    }
}
