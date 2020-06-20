package model

import (
	"github.com/gin-gonic/gin"
)

/*
Scan : Scan Object
*/
type Scan struct {
	Name  string                 `form:"scan" bson:"name"`
	Scope []string               `form:"scope[]" bson:"scope"`
	Modes map[string]interface{} `form:"modes[]" bson:"-"`
	// Result  map[string]interface{} `form:"results[]" bson:"result"`
	Result  []string
	Threads int               `bson:"threads"`
	Raw     string            `bson:"raw"`
	Options map[string]string `form:"options[]" bson:"options"`
}

func (s *Scan) LoadModes() {
	var funcMap = map[string]interface{}{
		"SubdomainScan": SubdomainScan,
		"HttProbe":      HttProbe,
		"DirsearchScan": DirsearchScan,
	}
	s.Modes = funcMap
}

func (s *Scan) GetName() string {
	return s.Name
}

// func (s *Scan) SetResult(res map[string]interface{}) {
func (s *Scan) SetResult(res []string) {

	s.Result = res
}

func (s *Scan) SetRaw(res string) {

	s.Raw = res
}

func (s *Scan) Run(args ...interface{}) (interface{}, error) {
	result, raw, err := s.Modes[s.Name].(func(s *Scan, args ...interface{}) ([]string, string, error))(s, args)

	if err != nil {
		return nil, err
	}
	s.Result = result
	s.Raw = raw
	return result, nil
}

func CreateBindScan(c *gin.Context, program *Program) (scan *Scan, err error) {
	c.Bind(&scan)
	options := c.PostFormMap("options")
	scan.Options = options
	scan.Threads = program.Threads
	scan.LoadModes()
	return scan, nil
}
