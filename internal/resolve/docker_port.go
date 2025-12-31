package resolve

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ResolveDockerByPort tries to find a container that publishes the given host port
func ResolveDockerByPort(port int) (*DockerOwner, bool, error) {
	cmd := exec.Command(
		"docker", "ps",
		"--format", "{{.ID}}",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, false, nil
	}

	for _, id := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if id == "" {
			continue
		}

		inspect := exec.Command("docker", "inspect", id)
		var data bytes.Buffer
		inspect.Stdout = &data
		inspect.Stderr = &data

		if err := inspect.Run(); err != nil {
			continue
		}

		txt := data.String()

		needle := fmt.Sprintf(`"HostPort": "%d"`, port)
		if !strings.Contains(txt, needle) {
			continue
		}

		name := extractJSONField(txt, `"Name":`)
		image := extractJSONField(txt, `"Image":`)

		composeProj := extractJSONLabel(txt, "com.docker.compose.project")
		composeSvc := extractJSONLabel(txt, "com.docker.compose.service")

		return &DockerOwner{
			ContainerID:   id[:12],
			ContainerName: strings.TrimPrefix(name, "/"),
			Image:         image,
			PID:           0,
			ComposeProj:   composeProj,
			ComposeSvc:    composeSvc,
		}, true, nil
	}

	return nil, false, nil
}

