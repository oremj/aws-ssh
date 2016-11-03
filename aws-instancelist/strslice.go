package main

type stringSliceVar []string

func (s *stringSliceVar) String() string {
	return ""
}

func (s *stringSliceVar) Set(val string) error {
	*s = append(*s, val)
	return nil
}
