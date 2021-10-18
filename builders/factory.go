package builders

func CreateBuilder(t string) interface{} {
	if t == "single_file" {
		return createSingleFileBuilder()
	}

	if t == "project" {
		return createProjectBuilder()
	}

	return nil
}

func CreateDestroyer() interface{} {
	return createSingleFileDestroyer()
}