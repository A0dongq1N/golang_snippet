package gormsnippet

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no .env file found, skip")
	}
}

// 定义模型
type User struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	Email     string
	CreatedAt time.Time
}

// 将User关联到t_user
func (User) TableName() string {
	return "t_user"
}

func TestGormQuery(t *testing.T) {
	// 连接数据库
	dsn := os.Getenv("GORM_SNIPPET_DSN")
	t.Logf("dsn=%v", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// 自动迁移 schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		return
	}

	// 创建记录
	user := User{Name: "张三"}
	db.Create(&user)

	// 查询单条记录
	var result User
	db.First(&result) // 获取第一条记录
	fmt.Printf("First record: %v\n", result)

	// 条件查询
	var users []User
	db.Where("id = ?", 5).Find(&users)
	fmt.Printf("Users id = 5: %v\n", users)

	// 更新记录
	db.Model(&result).Update("name", "李四")

	// 删除记录
	//db.Delete(&result)
}

func TestGormSelect(t *testing.T) {
	// 连接数据库 - 使用正确的主机名和端口
	// 相当于: mysql -hm1.idacdb-test.bus -uroot -proot1234 -P37256
	dsn := os.Getenv("GORM_SNIPPET_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 测试连接是否正常
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("获取底层数据库连接失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("数据库连接测试失败: %v", err)
	}

	t.Log("数据库连接成功")

	// 在SQL中明确指定数据库和表名
	rows, err := db.Raw("SELECT * FROM t_user WHERE id=?", 1).Rows()
	if err != nil {
		t.Fatalf("执行SQL查询失败: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("获取列信息失败: %v", err)
	}

	// 处理查询结果
	for rows.Next() {
		// 创建切片来存储行数据
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(scanArgs...); err != nil {
			t.Fatalf("扫描行数据失败: %v", err)
		}

		// 输出记录内容
		record := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// 处理字节数组转换为字符串
			if b, ok := val.([]byte); ok {
				record[col] = string(b)
			} else {
				record[col] = val
			}
		}

		t.Logf("记录内容: %+v", record)
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("遍历结果集时出错: %v", err)
	}
}
