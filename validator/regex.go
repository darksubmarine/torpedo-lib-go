package validator

import "regexp"

type Regex struct {
	re  *regexp.Regexp
	val string
}

func NewRegex(exp string) *Regex {
	return &Regex{re: regexp.MustCompile(exp)}
}

func (r *Regex) Value(val interface{}) IsValidInterface {
	if v, ok := val.(string); ok {
		r.val = v
	} else {
		r.val = ""
	}

	return r
}

func (r *Regex) IsValid() bool {
	return r.re.MatchString(r.val)
}
