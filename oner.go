// Imlementation of One_R

package godatamining

type OneR struct {
	Data       *[][]string
	MinBuckets int
}

type OneResult struct {
	Key        string
	TotalError float32
	Rules      *map[string]map[string]float32
}

func (oner *OneR) getKeyLoc(key string) (int, bool) {
	keys := (*oner.Data)[0]
	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			return i, true
		}
	}
	return -1, false
}

func (oner *OneR) GetInstance(key string, class string) (results map[string]map[string]int, ok bool) {
	loc, ok := oner.getKeyLoc(key)
	resultloc, ok := oner.getKeyLoc(class)
	// resultloc := len((*oner.Data)[0]) - 1
	if !ok || loc == resultloc {
		ok = false
		return
	}
	results = make(map[string]map[string]int)

	for i := 1; i < len(*oner.Data); i++ {
		if _, ok := results[(*oner.Data)[i][loc]]; !ok {
			results[(*oner.Data)[i][loc]] = map[string]int{
				(*oner.Data)[i][resultloc]: 1}
		} else {
			results[(*oner.Data)[i][loc]][(*oner.Data)[i][resultloc]] += 1
		}

	}
	return
}

func (oner *OneR) GetRules(key string, class string) (*map[string]string, bool) {
	if data, ok := oner.GetInstance(key, class); ok {
		rules := make(map[string]string)
		for a, b := range data {
			rules[a] = ""
			for x, y := range b {
				if _, ok := rules[a]; !ok {
					rules[a] = x
				} else {
					if b[rules[a]] < y {
						rules[a] = x
					}
				}
			}
		}
		// exclude rules where every value of a class is unique
		if len(rules) != len(*oner.Data)-1 {
			ok = true
			return &rules, true
		}
	}
	return new(map[string]string), false
}

func (oner *OneR) GetErrorRate(key string, class string) (*map[string]map[string]float32, bool) {
	if rules, ok := oner.GetRules(key, class); ok {
		output := make(map[string]map[string]float32)
		loc, _ := oner.getKeyLoc(key)
		n := len(*oner.Data)
		resultloc, _ := oner.getKeyLoc(class)
		// resultloc := len((*oner.Data)[0]) - 1
		for a, b := range *rules {
			mistakes := 0
			instances := 0
			for i := 1; i < n; i++ {
				if (*oner.Data)[i][loc] == a {
					instances += 1
					if (*oner.Data)[i][resultloc] != b {
						mistakes += 1
					}
				}
			}
			output[a] = map[string]float32{b: float32(mistakes) / float32(n-1)}
		}
		return &output, true
	}
	return new(map[string]map[string]float32), false
}

func (oner *OneR) GetBestRule(class string) (result *OneResult, ok bool) {
	output := make(map[string]*OneResult)
	for _, key := range (*oner.Data)[0] {
		if rules, ok := oner.GetErrorRate(key, class); ok {
			output[key] = &OneResult{key, float32(0), rules}
			for _, mistake := range *rules {
				for _, z := range mistake {
					output[key].TotalError += z
				}

			}
		}
	}

	for _, b := range output {
		if result == nil {
			result = b
		} else {
			if result.TotalError > b.TotalError {
				result = b
			}
		}
	}

	ok = true
	return

}
