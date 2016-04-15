package hoetom

const (
	UNKNOW = iota
	BOY
	GIRL
)

type Sex struct {
	Value int64
}

func SexFromString(cnstring string) Sex {
	switch cnstring {
	case "男":
		return Sex{BOY}
	case "女":
		return Sex{GIRL}
	default:
		return Sex{UNKNOW}
	}
}

func (this Sex) String() string {
    
	return this.CnString()
}

func (this Sex) CnString() string {
	if this.Value == BOY {
		return "男"
	} else if this.Value == GIRL {
		return "女"
	} else {
		return "未知"
	}
}
