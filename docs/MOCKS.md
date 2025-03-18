# Using Mockery for Testing

This project uses [mockery](https://github.com/vektra/mockery) to generate mock implementations of interfaces for testing.

## Generating Mocks

To generate mocks for all interfaces defined in the `.mockery.yaml` configuration file, run:

```bash
make mocks
```

This will generate mock implementations in the `internal/apps/contracts/mocks` directory.

## Configuration

The mockery configuration is defined in `.mockery.yaml` at the root of the project. This file specifies which interfaces to generate mocks for and where to place the generated files.

Current configuration targets:
- `FizzBuzzServiceIface`
- `StatsServiceIface`

## Using Mocks in Tests

Here's an example of how to use the generated mocks in your tests:

```go
import (
    "testing"
    "github.com/stretchr/testify/mock"
    "fizzbuzz-server/internal/apps/contracts/mocks"
)

func TestSomething(t *testing.T) {
    // Create a new mock instance
    mockFizzBuzzService := &mocks.FizzBuzzServiceIface{}
    
    // Set expectations
    mockFizzBuzzService.On("GenerateFizzBuzz", 3, 5, 15, "fizz", "buzz").
        Return([]string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz"})
    
    // Use the mock in your test
    // ...
    
    // Verify that all expectations were met
    mockFizzBuzzService.AssertExpectations(t)
}
```

## Adding New Interfaces

If you add new interfaces to the `contracts` package, you'll need to update the `.mockery.yaml` file to include them, then run `make mocks` again to generate the new mock implementations.