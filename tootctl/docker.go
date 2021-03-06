package tootctl

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"github.com/mohemohe/lxbot-mastoadmin/util"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	lines = 5
)

func deepCopy(msg util.M) (util.M, error) {
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	d := gob.NewDecoder(&b)
	if err := e.Encode(msg); err != nil {
		return nil, err
	}
	r := map[string]interface{}{}
	if err := d.Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func generateText(stdout []string, tag string) string {
	stdoutText := strings.Join(stdout, "\n")
	return tag + "\n\n" + stdoutText
}

func Run(msg util.M, script string, ch *chan util.M) {
	args := []string{"run", "--rm", "-i", "--log-driver", "none"}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "PATH") {
			continue
		}
		args = append(args, "-e", v)
	}
	image := os.Getenv("LXBOT_MASTOADMIN_DOCKER_IMAGE")
	if image == "" {
		image = "tootsuite/mastodon"
	}
	args = append(args, image)
	args = append(args, strings.Fields(script)...)

	log.Println("docker", strings.Join(args, " "))

	cmd := exec.Command("docker", args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	stdoutBuff := util.NewBuff()
	buffCh := make(chan int)
	timeout := false

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdoutBuff.Enqueue(scanner.Text())
			buffCh <- 1
		}

		wg.Done()
	}()
	go func() {
		wg.Add(1)

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stdoutBuff.Enqueue(scanner.Text())
			buffCh <- 1
		}

		wg.Done()
	}()
	go func() {
		waitCh := time.After(3 * time.Second)

		for {
			select {
			case _, ok := <-buffCh:
				if !ok {
					return
				}
				stdoutBuffLen := stdoutBuff.Len()
				if stdoutBuffLen >= lines {
					waitCh = time.After(3 * time.Second)

					stdoutLen := lines
					if stdoutBuffLen < lines {
						stdoutLen = stdoutBuffLen
					}

					stdoutLines := stdoutBuff.BulkDequeue(stdoutLen)
					text := generateText(stdoutLines, "(PARTIAL)")

					// FIXME: copy error
					nextMsg, _ := deepCopy(msg)
					nextMsg["mode"] = "reply"
					nextMsg["message"].(util.M)["text"] = text
					*ch <- nextMsg
				}
				break
			case <-waitCh:
				waitCh = time.After(3 * time.Second)

				stdoutLines := stdoutBuff.DequeueALL()
				if len(stdoutLines) == 0 {
					break
				}
				text := generateText(stdoutLines, "(PARTIAL)")

				// FIXME: copy error
				nextMsg, _ := deepCopy(msg)
				nextMsg["mode"] = "reply"
				nextMsg["message"].(util.M)["text"] = text
				*ch <- nextMsg
				break
			case <-time.After(120 * time.Minute):
				if !cmd.ProcessState.Exited() {
					_ = cmd.Process.Kill()
				}
				timeout = true
				return
			}
		}
	}()

	wg.Wait()
	_ = cmd.Wait()
	close(buffCh)

	tag := "(FINISH)"
	if timeout {
		tag = "(TIMEOUT)"
	}
	stdoutLines := stdoutBuff.DequeueALL()
	text := generateText(stdoutLines, tag)

	// FIXME: copy error
	nextMsg, err := deepCopy(msg)
	if err != nil {
		log.Println(err)
		return
	}
	nextMsg["mode"] = "reply"
	nextMsg["message"].(util.M)["text"] = text
	*ch <- nextMsg
}
