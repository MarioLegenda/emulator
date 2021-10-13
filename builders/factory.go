package builders

func CreateBuilder(t string) interface{} {
	if t == "blog" {
		return createSingleFileBuilder()
	}
	return nil
}