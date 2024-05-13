package utility

func SplitEveryN(s string, interval int) []string {
        substrings := []string{}
        for i := 0; i < len(s); i++ {
                if i % interval == 0 {
                        substrings = append(substrings, "")
                }

                substrings[len(substrings) - 1] += string(s[i])
        }
        return substrings
}

