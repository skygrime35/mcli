// internal/docker/list_test.go
package docker

import "testing"

const sampleDockerPsOutput = `7d35e839c844|eva_preprod_frontend|Up 3 hours (unhealthy)|eva_preprod-frontend
10637e7259e7|eva_preprod_backend|Up 3 hours (healthy)|eva_preprod-backend
a33f361fb87e|eva_preprod_redis|Up 3 hours (healthy)|redis:7-alpine
`

func TestParseContainerList(t *testing.T) {
	containers := parseContainerList(sampleDockerPsOutput)
	if len(containers) != 3 {
		t.Fatalf("expected 3 containers, got %d", len(containers))
	}
	if containers[0].ID != "7d35e839c844" || containers[0].Name != "eva_preprod_frontend" ||
		containers[0].Status != "Up 3 hours (unhealthy)" || containers[0].Image != "eva_preprod-frontend" {
		t.Errorf("container 0 parsed incorrectly: %+v", containers[0])
	}
	if containers[2].Image != "redis:7-alpine" {
		t.Errorf("container 2 image = %q, want redis:7-alpine", containers[2].Image)
	}
}

func TestParseContainerList_Empty(t *testing.T) {
	if got := parseContainerList(""); len(got) != 0 {
		t.Errorf("expected 0 containers for empty input, got %d", len(got))
	}
}

func TestParseContainerList_MalformedLineSkipped(t *testing.T) {
	containers := parseContainerList("not-enough-fields\n" + sampleDockerPsOutput)
	if len(containers) != 3 {
		t.Errorf("expected malformed line to be skipped, got %d containers", len(containers))
	}
}
