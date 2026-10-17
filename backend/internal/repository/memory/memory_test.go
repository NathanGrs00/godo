package memory_test

import (
	"testing"

	"github.com/NathanGrs00/godo/internal/repository/memory"
	"github.com/NathanGrs00/godo/internal/repository/repotest"
)

func TestContract(t *testing.T) {
	repotest.Run(t, func(*testing.T) repotest.Repositories {
		return repotest.Repositories{
			Tasks:     memory.NewTaskRepository(),
			Deadlines: memory.NewDeadlineRepository(),
			Tags:      memory.NewTagRepository(),
			Users:     memory.NewUserRepository(),
		}
	})
}
