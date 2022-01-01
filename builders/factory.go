package builders

func CreateBuilder(t string) interface{} {
	if t == "single_file" {
		return createSingleFileBuilder()
	}

	if t == "project" {
		return createProjectBuilder()
	}

	if t == "c_project" {
		return createCLangBuilder()
	}

	if t == "linked_compiled_project" {
		return createCompiledProject()
	}

	return nil
}

func CreateDestroyer(t string) interface{} {
	if t == "single_file" {
		return createSingleFileDestroyer()
	}

	if t == "project" {
		return createProjectDestroyer()
	}

	return nil
}
