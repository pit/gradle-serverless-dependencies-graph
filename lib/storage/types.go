package storage

type (
	DependenciesRest struct {
		Dependencies []DependencyRest `json:"dependencies"`
	}

	DependencyRest struct {
		Group   string `json:"group"`
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	UpsertResultRest struct {
		UsedCapacity float64
	}

	StorageErrorRest struct {
		Message string
		Code    int
		Repo    string
		Ref     string
		Id      string
		Version string
		Err     error
	}
)

type (
	DependencyDto struct {
		Parent string `dynamodbav:"Parent"`
		Child  string `dynamodbav:"Child"`
	}

	RepositoryDto struct {
		Parent string `dynamodbav:"Parent"`
		Child  string `dynamodbav:"Child"`
	}

	StorageDto struct {
		Dependency string `dynamodbav:"Dependency"`
		Version    string `dynamodbav:"Version"`
		Repo       string `dynamodbav:"Repo"`
		Ref        string `dynamodbav:"Ref"`
	}
)
