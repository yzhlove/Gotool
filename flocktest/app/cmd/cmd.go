package cmd

type Arg struct {
	Key  string
	Vars []string
}

type Args []Arg

func (a Args) V() []string {
	res := make([]string, 0, len(a)*2)
	for _, t := range a {
		res = append(res, t.Key)
		if len(t.Vars) != 0 {
			res = append(res, t.Vars...)
		}
	}
	return res
}

func Wrap(trail string, values ...Arg) Args {
	values = append(values, Arg{Key: trail})
	return values
}
