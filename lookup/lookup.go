package lookup

import (
	"blogpost/models"
	migrators "blogpost/updates"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func LookUp(db *migrators.LookUpDb) {
	update := migrators.LookUpDb{
		DB: db.DB,
	}

	folderPath := "./updates"
	lookUp := &models.LookUp{}

	//Get a list of file info objects in the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Count the number of files in the folder, excluding files with the name "lookup1"

	fileCount := 0
	filesName := make([]string, 0)

	for _, file := range files {
		if file.IsDir() || file.Name() == "migrations.go" {
			// Skip directories and files with the name "lookup1"
			continue
		}
		fileCount++
		filesName = append(filesName, file.Name())
	}

	fmt.Println("filesName", filesName)

	for _, file := range filesName {
		fmt.Println("file name", file)
		version, err := extractNumberFromFileName(file)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if err := db.DB.Where("name", file).Where("version", version).First(&lookUp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.DB.Create(&models.LookUp{
					Name:    file,
					Version: version,
				}).Error; err != nil {
					if errors.Is(err, gorm.ErrDuplicatedKey) {
						continue
					}
				}

				fmt.Println("Andser-->", strings.Title(file))
				extractedName := extractNameFromFileName(file)
				methodName := strings.Title(extractedName)
				fmt.Println("method name", methodName)

				method := reflect.ValueOf(&update).MethodByName(methodName)
				if !method.IsValid() {
					fmt.Printf("Method %s not found or not exported\n", methodName)
					continue
				}

				// Check if the receiver is a pointer and if it's initialized
				if method.Type().NumIn() != 0 && method.Type().In(0).Kind() == reflect.Ptr {
					// Create an instance of the receiver type
					receiverType := method.Type().In(0).Elem()
					receiverInstance := reflect.New(receiverType).Elem()

					// Check if it's nil
					if receiverInstance.Interface() == nil {
						fmt.Println("Receiver is not properly initialized")
						continue
					}
				}

				// Call the method
				method.Call([]reflect.Value{})
			}
		}
	}
}

func extractNumberFromFileName(fileName string) (int, error) {
	// Use a regular expression to find the number in the filename
	re := regexp.MustCompile(`(\d+)`)
	match := re.FindStringSubmatch(fileName)

	if len(match) < 2 {
		return 0, fmt.Errorf("number not found in filename")
	}

	// Convert the matched string to an integer
	number, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}

	return number, nil
}

func extractNameFromFileName(fileName string) string {
	// Use filepath.Base to extract the base name of the file without the extension
	name := filepath.Base(fileName)
	// Remove the ".go" extension if present
	name = name[:len(name)-len(filepath.Ext(name))]
	return name
}
