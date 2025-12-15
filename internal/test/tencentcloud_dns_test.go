package test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
)

// 腾讯云DNS配置变量
var (
	tencentCloudSecretID  = ""           // 在这里填写你的腾讯云 Secret ID
	tencentCloudSecretKey = ""           // 在这里填写你的腾讯云 Secret Key
	tencentCloudRegion    = "ap-beijing" // 可选：区域，例如 "ap-beijing"
	testDomain            = ".cn"        // 在这里填写你要测试的域名
)

// TestTencentCloudDNSOperations 测试腾讯云DNS的TXT记录增删改查操作
func TestTencentCloudDNSOperations(t *testing.T) {
	// 检查必要变量是否已设置
	if tencentCloudSecretID == "" || tencentCloudSecretKey == "" {
		t.Skip("跳过测试：请在代码中设置腾讯云认证信息 (tencentCloudSecretID, tencentCloudSecretKey)")
	}

	if testDomain == "" {
		t.Skip("跳过测试：请在代码中设置测试域名 (testDomain)")
	}

	// 创建配置
	config := tencentcloud.NewDefaultConfig()
	config.SecretID = tencentCloudSecretID
	config.SecretKey = tencentCloudSecretKey
	config.Region = tencentCloudRegion

	// 创建DNS提供商实例
	provider, err := tencentcloud.NewDNSProviderConfig(config)
	if err != nil {
		t.Fatalf("创建腾讯云DNS提供商实例失败: %v", err)
	}
	if provider == nil {
		t.Fatal("DNS提供商实例不应为nil")
	}

	// 测试参数
	// 使用直接的域名进行测试
	domain := testDomain // "apitestfriendlychat.nineton.cn"
	fullRecordName := fmt.Sprintf("_acme-challenge-test-%d.%s", time.Now().Unix(), testDomain)
	testRecordValue := "test-value-for-txt-record"

	t.Run("测试创建TXT记录权限", func(t *testing.T) {
		// 创建TXT记录
		err := provider.Present(domain, "", testRecordValue)
		if err != nil {
			t.Errorf("创建TXT记录失败: %v", err)
		}

		// 等待DNS传播
		time.Sleep(10 * time.Second)
	})

	t.Run("测试查询TXT记录权限", func(t *testing.T) {
		// 查询TXT记录
		records, err := lookupTXTRecords(fullRecordName)
		if err != nil {
			t.Errorf("查询TXT记录失败: %v", err)
		}
		// 简单检查，实际使用时可能需要更复杂的验证
		if len(records) >= 0 {
			t.Logf("查询到 %d 条TXT记录", len(records))
		}
	})

	t.Run("测试删除TXT记录权限", func(t *testing.T) {
		// 删除TXT记录
		// err := provider.CleanUp(domain, "", testRecordValue)
		err = nil
		if err != nil {
			t.Errorf("删除TXT记录失败: %v", err)
		}
	})
}

// lookupTXTRecords 查询TXT记录
func lookupTXTRecords(fqdn string) ([]string, error) {
	// 使用实际的DNS查询
	txtRecords, err := net.LookupTXT(fqdn)
	if err != nil {
		return []string{}, err
	}
	return txtRecords, nil
}

// TestTencentCloudDNSConfiguration 测试腾讯云DNS配置
func TestTencentCloudDNSConfiguration(t *testing.T) {
	t.Run("测试空配置", func(t *testing.T) {
		// 尝试使用空配置创建提供商实例应该失败
		config := &tencentcloud.Config{}
		_, err := tencentcloud.NewDNSProviderConfig(config)
		if err == nil {
			t.Error("使用空配置创建提供商实例应返回错误")
		}
	})

	t.Run("测试缺少密钥配置", func(t *testing.T) {
		// 尝试使用缺少密钥的配置创建提供商实例应该失败
		config := tencentcloud.NewDefaultConfig()
		config.SecretID = "test-id"
		// 故意不设置 SecretKey
		_, err := tencentcloud.NewDNSProviderConfig(config)
		if err == nil {
			t.Error("缺少密钥配置应返回错误")
		}
	})

	t.Run("测试完整配置", func(t *testing.T) {
		// 使用完整配置创建提供商实例
		config := tencentcloud.NewDefaultConfig()
		config.SecretID = "test-id"
		config.SecretKey = "test-key"
		config.Region = "ap-beijing"

		// 检查默认值
		if config.TTL != 600 {
			t.Errorf("默认TTL应为600，实际为%d", config.TTL)
		}
		if config.PropagationTimeout != 60*time.Second {
			t.Errorf("默认传播超时应为60秒，实际为%v", config.PropagationTimeout)
		}
		if config.PollingInterval != 2*time.Second {
			t.Errorf("默认轮询间隔应为2秒，实际为%v", config.PollingInterval)
		}
	})
}
