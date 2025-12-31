package resolve

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ResolveDockerByPort tries to find a container that publishes the given host port
func ResolveDockerByPort(port int) (*DockerOwner, bool, error) {
	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		// docker not available or not running => not resolvable here
		return nil, false, nil
	}

	needle := fmt.Sprintf(`"HostPort": "%d"`, port)

	for _, id := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if id == "" {
			continue
		}

		txt, err := dockerInspect(id)
		if err != nil {
			continue
		}

		if !strings.Contains(txt, needle) {
			continue
		}

		name := extractJSONField(txt, `"Name":`)
		image := extractJSONField(txt, `"Image":`)
		composeProj := extractJSONLabel(txt, "com.docker.compose.project")
		composeSvc := extractJSONLabel(txt, "com.docker.compose.service")

		return &DockerOwner{
			ContainerID:   id, // FULL
			ContainerName: strings.TrimPrefix(name, "/"),
			Image:         image,
			PID:           0, // host-side proxy/dockerd, not container PID
			ComposeProj:   composeProj,
			ComposeSvc:    composeSvc,
		}, true, nil
	}

	return nil, false, nil
}

