package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = "todo_test.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")

	os.Setenv("TODO_FILENAME", fileName)
	defer os.Unsetenv("TODO_FILENAME")

	result := m.Run()

	fmt.Println("Clearning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	task := "test task number 1"
	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	task2 := "test task number 2"
	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdin, task2)
		cmdStdin.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	task3 := "test task number 3"
	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task3)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	task4 := "test task number 4"
	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task4)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(cmdPath, "-del", "3")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")

		out, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("X 1: %s\n  2: %s\n  3: %s\n", task, task2, task4)

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead", expected, string(out))
		}
	})
}
