/*
Copyright the Sonobuoy contributors 2019

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package docker

import (
	"github.com/heptio/sonobuoy/pkg/image/exec"
	log "github.com/sirupsen/logrus"
)

type Docker interface {
	PullIfNotPresent(image string, retries int) error
	Pull(image string, retries int) error
	Push(image string, retries int) error
	Tag(src, dest string, retries int) error
	Rmi(image string, retries int) error
	Save(images []string, filename string) error
}

type LocalDocker struct {
}

// PullIfNotPresent will pull an image if it is not present locally
// retrying up to retries times
// returns errors from pulling
func (l LocalDocker) PullIfNotPresent(image string, retries int) error {
	// TODO(bentheelder): switch most (all) of the logging here to debug level
	// once we have configurable log levels
	// if this did not return an error, then the image exists locally
	cmd := exec.Command("docker", "inspect", "--type=image", image)
	if err := cmd.Run(); err == nil {
		log.Infof("Image: %s present locally", image)
		return nil
	}
	// otherwise try to pull it
	return l.Pull(image, retries)
}

// Pull pulls an image, retrying up to retries times
func (l LocalDocker) Pull(image string, retries int) error {
	log.Infof("Pulling image: %s ...", image)
	return exec.RunLoggingOutputOnFail(exec.Command("docker", "pull", image), retries)
}

// Push pushes an image, retrying up to retries times
func (l LocalDocker) Push(image string, retries int) error {
	log.Infof("Pushing image: %s ...", image)
	return exec.RunLoggingOutputOnFail(exec.Command("docker", "push", image), retries)
}

// Tag tags an image, retrying up to retries times
func (l LocalDocker) Tag(src, dest string, retries int) error {
	log.Infof("Tagging image: %s as %s ...", src, dest)
	return exec.RunLoggingOutputOnFail(exec.Command("docker", "tag", src, dest), retries)
}

// Rmi removes an image, retrying up to retries times
func (l LocalDocker) Rmi(image string, retries int) error {
	log.Infof("Deleting image: %s ...", image)
	return exec.RunLoggingOutputOnFail(exec.Command("docker", "rmi", image), retries)
}

// Save exports a set of images to a tar file
func (l LocalDocker) Save(images []string, filename string) error {
	log.Info("Saving images: ...")

	//TODO(stevesloka) Check if all images exist on local client first

	// Build out docker command
	args := append([]string{"save"}, images...)
	args = append(args, "--output", filename)

	return exec.RunLoggingOutputOnFail(exec.Command("docker", args...), 0)
}
