package cmds

type Arg struct {
	Key string
	Var []string
}

type Args []Arg

type Cmd struct {
	Params string
	Args   Args
}

func (c Cmd) V() []string {
	resp := make([]string, 0, len(c.Args)*2+1)
	for _, arg := range c.Args {
		resp = append(resp, arg.Key)
		resp = append(resp, arg.Var...)
	}
	if len(c.Params) != 0 {
		resp = append(resp, c.Params)
	}
	return resp
}

type Cmds []Cmd

func (c *Cmds) Add(cc ...Cmd) {
	*c = append(*c, cc...)
}

func (c *Cmds) V() []string {
	resp := make([]string, 0, len(*c)*2)
	for _, arg := range *c {
		resp = append(resp, arg.V()...)
	}
	return resp
}

func Wrap(command string, args ...Arg) Cmd {
	c := Cmd{
		Params: command,
	}
	for _, arg := range args {
		c.Args = append(c.Args, Arg{Key: arg.Key, Var: arg.Var})
	}
	return c
}
