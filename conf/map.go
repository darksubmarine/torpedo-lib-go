package conf

import "fmt"

type Map map[string]interface{}

func (c Map) FetchSubMapP(key ...string) Map {
	if val, ok := c.FetchP(key...).(Map); ok {
		return val
	} else {
		panic(fmt.Sprintf("invalid data type"))
	}
}

func (c Map) FetchP(key ...string) interface{} {
	if v, ok := c.Fetch(key...); !ok {
		panic(fmt.Sprintf("key not found: %v", key))
	} else {
		return v
	}
}

func (c Map) FetchStringP(key ...string) string {
	if val, ok := c.FetchP(key...).(string); ok {
		return val
	} else {
		panic(fmt.Sprintf("invalid data type"))
	}
}

func (c Map) FetchStringOrElse(val string, key ...string) string {
	if rawv, ok := c.Fetch(key...); ok {
		if v, ok := rawv.(string); ok {
			return v
		}
	}
	return val
}

func (c Map) FetchIntP(key ...string) int {
	if val, ok := c.FetchP(key...).(int); ok {
		return val
	} else {
		panic(fmt.Sprintf("invalid data type"))
	}
}

func (c Map) FetchIntOrElse(val int, key ...string) int {
	if rawv, ok := c.Fetch(key...); ok {
		if v, ok := rawv.(int); ok {
			return v
		}
	}
	return val
}

func (c Map) FetchBoolP(key ...string) bool {
	if val, ok := c.FetchP(key...).(bool); ok {
		return val
	} else {
		panic(fmt.Sprintf("invalid data type"))
	}
}

func (c Map) FetchBoolOrElse(val bool, key ...string) bool {
	if rawv, ok := c.Fetch(key...); ok {
		if v, ok := rawv.(bool); ok {
			return v
		}
	}
	return val
}

func (c Map) Add(val interface{}, key ...string) {
	prev := c
	for i, p := range key {
		if v, ok := prev[p]; ok {
			switch vv := v.(type) {
			case Map:
				if i == len(key)-1 {
					prev[p] = val
					//return vv, true
				}
				prev = vv
			case map[string]interface{}:
				if i == len(key)-1 {
					prev[p] = val
					//return vv, true
				}
				prev = vv

			default:
				prev[p] = val
				//return vv, true
			}
		} else {
			if i == len(key)-1 {
				prev[p] = val
				//return vv, true
			} else {
				prev[p] = Map{}
				prev = prev[p].(Map)
			}
		}
	}
	//return nil, false
}

func (c Map) Fetch(key ...string) (interface{}, bool) {
	prev := c
	for i, p := range key {
		if v, ok := prev[p]; ok {
			switch vv := v.(type) {
			case Map:
				if i == len(key)-1 {
					return vv, true
				}
				prev = vv
			case map[string]interface{}:
				if i == len(key)-1 {
					return vv, true
				}
				prev = vv

			default:
				return vv, true
			}
		}
	}
	return nil, false
}

func (c Map) Interpolate(verbose bool) Map {

	for k, v := range c {
		switch vt := v.(type) {
		case Map:
			if verbose {
				fmt.Printf("---%s\n", k)
			}
			vt.Interpolate(verbose)
		case string:
			vv := processWithEnvVarInterpolation(vt)
			if verbose {
				if vt != vv {
					fmt.Printf("%s=%v from %s\n", k, vv, v)
				} else {
					fmt.Printf("%s=%v\n", k, vv)
				}

			}
			c[k] = vv
		default:
			if verbose {
				fmt.Printf("%s=%v\n", k, v)
			}
		}
	}

	return c
}
