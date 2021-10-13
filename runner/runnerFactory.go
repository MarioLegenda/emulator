package runner

func CreateRunner(t string) interface{} {
	if t == "singleFile" {
		return createSingleFileRunner()
	}

	return nil
}
