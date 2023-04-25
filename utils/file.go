package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetAllFiles(dir string) ([]string, error) {
	var files []string
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, dirEntry := range dirEntries {
		fileName := dirEntry.Name()
		filePath := filepath.Join(dir, fileName)
		if dirEntry.IsDir() {
			subFiles, err := GetAllFiles(filePath)
			if err != nil {
				return files, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, filePath)
		}
	}
	return files, nil
}

func SaveProperties(properties map[string]string) error {
	file, err := os.Create(GetConfigPath())
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for key, value := range properties {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadProperties(path string) (map[string]string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create properties file: %v", err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		_, err = fmt.Fprintln(writer, "user.name="+GetOsUserName())
		if err != nil {
			return nil, fmt.Errorf("failed to write properties to file: %v", err)
		}
		_, err = fmt.Fprintln(writer, "auto.receive=false")
		if err != nil {
			return nil, fmt.Errorf("failed to write properties to file: %v", err)
		}
		err = writer.Flush()
		if err != nil {
			return nil, fmt.Errorf("failed to flush writer: %v", err)
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open properties file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	properties := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		properties[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read properties file: %v", err)
	}

	return properties, nil
}
func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return fileInfo.IsDir()
}
func CreateDirectories(path string) error {
	var dir = ""
	if IsDirectory(path) {
		dir = path
	} else {
		dir = filepath.Dir(path)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directories for file %s: %v", path, err)
		}
	}
	return nil
}
func GetExecutableDir() (string, error) {
	exePath, _ := os.Executable()
	return filepath.Dir(exePath), nil
}
