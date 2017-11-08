package awsutils

type StringSliceVar []string

func (s *StringSliceVar) String() string {
	return ""
}

func (s *StringSliceVar) Set(val string) error {
	*s = append(*s, val)
	return nil
}
