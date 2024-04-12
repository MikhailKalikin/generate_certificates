package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Выбор типа сертификата
	fmt.Println("Выберите тип сертификата:")
	fmt.Println("1. Клиентский")
	fmt.Println("2. Серверный")
	fmt.Println("3. Клиентский и серверный")
	fmt.Print("Введите номер: ")
	choice, _ := reader.ReadString('\n')

	// Ввод CN
	fmt.Print("Введите CN: ")
	cn, _ := reader.ReadString('\n')
	cn = cn[:len(cn)-1] // Удаление символа новой строки

	// Путь к текущему каталогу
	currentDir, _ := os.Getwd()

	// Путь для сохранения ключей и CSR
	keyPath := filepath.Join(currentDir, cn + "_keys")
	os.MkdirAll(keyPath, os.ModePerm)

	// Генерация ключа и CSR в зависимости от выбора
	switch choice {
	case "1\n":
		generateCSR("clientAuth", cn, keyPath)
	case "2\n":
		generateCSR("serverAuth", cn, keyPath)
	case "3\n":
		generateCSR("clientAuth,serverAuth", cn, keyPath)
	}
}

func generateCSR(usage, cn, path string) {
	privateKeyPath := filepath.Join(path, cn+".key")
	csrPath := filepath.Join(path, cn+".csr")

	// Генерация приватного ключа
	cmd := exec.Command("openssl", "genrsa", "-out", privateKeyPath, "2048")
	cmd.Run()

	// Создание CSR
	cmd = exec.Command("openssl", "req", "-new", "-key", privateKeyPath, "-out", csrPath, "-sha256", "-nodes",
		"-subj", fmt.Sprintf("/CN=%s", cn), "-addext", fmt.Sprintf("extendedKeyUsage=%s", usage))
	cmd.Run()

	fmt.Printf("Созданы ключ и CSR для %s, сохранены в %s\n", cn, path)
}
