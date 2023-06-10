// TODO: add negative test cases
package DI_test

import (
	"testing"
	"web-crawler/modules/DI"
)

type MockDependency struct {
	Name string
}

type MockService struct {
	Dependency  *MockDependency
	Dependency2 *MockDependency
}

var mockService *MockService

func beforeTest() {
	mockService = &MockService{
		Dependency: &MockDependency{
			Name: "Dependency1",
		},
		Dependency2: &MockDependency{
			Name: "Dependency2",
		},
	}
}

func TestRegisterInjectShouldWork(t *testing.T) {
	beforeTest()

	DI.Register(mockService, nil)

	mockService2 := &MockService{}
	DI.Inject(mockService2)

	if mockService2 == nil {
		t.Error("Injection failed: mockService2 is nil")
		return
	}

	if mockService2.Dependency == nil || mockService2.Dependency2 == nil {
		t.Error("Injection failed: dependencies are nil")
		return
	}

	if mockService2.Dependency.Name != mockService.Dependency.Name {
		t.Error("Injection failed: mockService2.Dependency.Name is not equal to mockService.Dependency.Name")
		return
	}

	if mockService2.Dependency2.Name != mockService.Dependency2.Name {
		t.Error("Injection failed: mockService2.Dependency2.Name is not equal to mockService.Dependency2.Name")
		return
	}
}
