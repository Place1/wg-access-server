package authsession

type claim struct {
	Name  string
	Value string
}

type Claims []claim

func (c *Claims) Add(name string, value string) {
	*c = append(*c, claim{
		Name:  name,
		Value: value,
	})
}

func (c *Claims) Contains(claim string) bool {
	for _, curr := range *c {
		if curr.Name == claim {
			return true
		}
	}
	return false
}

func (c *Claims) Has(claim string, value string) bool {
	for _, curr := range *c {
		if curr.Name == claim {
			if curr.Value == value {
				return true
			}
		}
	}
	return false
}
