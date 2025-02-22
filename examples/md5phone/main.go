package main

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	mu          sync.Mutex
	totalHashes int32
	matched     int32
	memoryMode  bool
)

// 定义一个用于存储数据的结构体
type PhoneHash struct {
	ID    int
	Phone string
	Hash  string
}

var phoneHashesInMemory []PhoneHash // 在内存模式下使用的切片

// 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	// 使用 GORM 打开数据库连接
	db, err := gorm.Open(sqlite.Open("phone_hashes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 从 GORM DB 对象获取原生数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动创建表
	db.AutoMigrate(&PhoneHash{})

	return db, nil
}

func main() {
	start := time.Now()

	// 定义命令行参数
	memoryFlag := flag.Bool("memory", true, "Use memory mode to store data instead of SQLite database")
	hashFilePath := flag.String("hashfile", "hashes.txt", "Path to the file containing MD5 hashes")
	prefixFilePath := flag.String("pre", "pre.txt", "Path to the file containing phone number prefixes")
	flag.Parse() // 解析命令行参数
	fmt.Println("默认为内存模式, 如果需要翻译的md5值数据量过大, 请使用 -memory=false 关闭内存模式, 将会使用sqlite3储存中间数据。")
	memoryMode = *memoryFlag
	// 读取哈希文件内容
	hashes, err := loadHashes(*hashFilePath)
	if err != nil {
		fmt.Println("Error loading hashes:", err)
		return
	}
	totalHashes = int32(len(hashes))

	// 启动一个协程定期打印进度
	go startProgress()

	// 读取手机号前缀文件
	prefixes, err := loadPrefixes(*prefixFilePath)
	if err != nil {
		fmt.Println("Error loading prefixes:", err)
		return
	}

	if memoryMode {
		runInMemoryMode(hashes, prefixes)
	} else {
		runInDBMode(hashes, prefixes)
	}
	duration := time.Since(start)
	fmt.Printf("\nProgram finished in %s\n", duration)
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func insertPhoneHash(db *gorm.DB, phone, hash string) {
	db.Create(&PhoneHash{Phone: phone, Hash: hash})
}

func loadHashes(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hashes := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hash := scanner.Text()
		hashes[hash] = struct{}{} // Use an empty struct to minimize memory usage
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hashes, nil
}

func loadPrefixes(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var prefixes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		prefix := strings.TrimSpace(scanner.Text())
		if prefix != "" {
			prefixes = append(prefixes, prefix)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return prefixes, nil
}

func exportToCSV(db *gorm.DB, csvFilePath string) error {
	var phoneHashes []PhoneHash
	result := db.Find(&phoneHashes)
	if result.Error != nil {
		return result.Error
	}

	file, err := os.Create(csvFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, ph := range phoneHashes {
		if err := writer.Write([]string{ph.Phone, ph.Hash}); err != nil {
			return err
		}
	}

	return nil
}

func writeUnmatchedHashes(hashes map[string]struct{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for hash := range hashes {
		if _, err := file.WriteString(hash + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// 将内存数据写入到文本文件
func writeToTextFile(data []PhoneHash, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, ph := range data {
		if _, err := file.WriteString(fmt.Sprintf("%s,%s\n", ph.Phone, ph.Hash)); err != nil {
			return err
		}
	}

	return nil
}

func runInMemoryMode(hashes map[string]struct{}, prefixes []string) {
	var wg sync.WaitGroup
	for _, prefix := range prefixes {
		wg.Add(1)
		go func(prefix string) {
			defer wg.Done()
			for i := 0; i < 100000000; i++ {
				phone := fmt.Sprintf("%s%08d", prefix, i)
				hash := md5Hash(phone)
				mu.Lock()
				if _, found := hashes[hash]; found {
					phoneHashesInMemory = append(phoneHashesInMemory, PhoneHash{Phone: phone, Hash: hash})
					atomic.AddInt32(&matched, 1)
					delete(hashes, hash)
				}
				mu.Unlock()
			}
		}(prefix)
	}
	wg.Wait()

	// ... 将内存中的数据写入文件 ...
	if err := writeToTextFile(phoneHashesInMemory, "output.txt"); err != nil {
		fmt.Println("Error writing to text file:", err)
	}
	// ... 将未匹配的哈希写入文件 ...
	if err := writeUnmatchedHashes(hashes, "unmatched.txt"); err != nil {
		fmt.Println("Error writing unmatched hashes:", err)
	}
}

func runInDBMode(hashes map[string]struct{}, prefixes []string) {
	db, err := InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	var wg sync.WaitGroup
	for _, prefix := range prefixes {
		wg.Add(1)
		go func(prefix string) {
			defer wg.Done()
			for i := 0; i < 100000000; i++ {
				phone := fmt.Sprintf("%s%08d", prefix, i)
				hash := md5Hash(phone)
				mu.Lock()
				if _, found := hashes[hash]; found {
					insertPhoneHash(db, phone, hash)
					atomic.AddInt32(&matched, 1)
					delete(hashes, hash)
				}
				mu.Unlock()
			}
		}(prefix)
	}
	wg.Wait()

	// 导出到 CSV 文件
	if err := exportToCSV(db, "output.csv"); err != nil {
		fmt.Println("Error exporting to CSV:", err)
		return
	}
	// 关闭数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Error closing database:", err)
		return
	}
	sqlDB.Close()

	// 删除数据库文件
	if err := os.Remove("phone_hashes.db"); err != nil {
		fmt.Println("Error deleting database file:", err)
	}
	// ... 将未匹配的哈希写入文件 ...
	if err := writeUnmatchedHashes(hashes, "unmatched.txt"); err != nil {
		fmt.Println("Error writing unmatched hashes:", err)
	}
}

func startProgress() {
	for {
		time.Sleep(1 * time.Second) // 每秒更新一次
		currentMatched := atomic.LoadInt32(&matched)
		fmt.Printf("\r已翻译: %d/%d", currentMatched, totalHashes)
		if currentMatched >= totalHashes {
			break
		}
	}
}
