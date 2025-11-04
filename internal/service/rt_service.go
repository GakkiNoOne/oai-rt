package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"rt-manage/internal/config"
	"rt-manage/internal/model"
	"rt-manage/internal/repository"
	"rt-manage/pkg/logger"

	http2 "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/google/uuid"
	"golang.org/x/net/proxy"
)

// RTService RT 服务接口
type RTService interface {
	Create(rt *model.RT) error
	Update(id int64, updates map[string]interface{}) (*model.RT, error)
	GetByID(id int64) (*model.RT, error)
	GetByBizId(bizId string) (*model.RT, error)
	GetByEmail(email string) (*model.RT, error)
	List(page, pageSize int, name string, tag string, email string, typeStr string, enabled *bool, createDate string) ([]*model.RT, int64, error)
	Delete(id int64) error
	BatchDelete(ids []int64) (int, int, error)
	Refresh(id int64, refreshUserInfo, refreshAccountInfo bool) (*model.RT, error)
	RefreshUserInfo(id int64) (*model.RT, error)
	RefreshAccountInfo(id int64) (*model.RT, error)
	BatchRefresh(ids []int64) (int, int, []map[string]interface{}, error)
	BatchImport(batchName string, tag string, tokens []string, proxyList []string) (int, int, error)
	AutoRefreshAll() error
}

type rtService struct {
	repo       repository.RTRepository
	configRepo repository.ConfigRepository
}

// NewRTService 创建 RT 服务实例
func NewRTService(repo repository.RTRepository, configRepo repository.ConfigRepository) RTService {
	return &rtService{
		repo:       repo,
		configRepo: configRepo,
	}
}

// createHTTPClient 创建支持 HTTP/HTTPS/SOCKS5 代理的 HTTP 客户端
// 配置为模拟浏览器行为，避免被 Cloudflare 识别
func createHTTPClient(proxyURL string, timeout time.Duration) (*http.Client, error) {
	// 基础 Transport 配置（模拟浏览器）
	transport := &http.Transport{
		// 禁用 HTTP/2，使用 HTTP/1.1（与 Python requests 行为一致）
		ForceAttemptHTTP2:     false,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		// 连接池配置
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     0,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		// 禁用自动压缩，我们手动设置 Accept-Encoding 以完全匹配 Python requests
		DisableCompression:  true,
	}

	if proxyURL != "" {
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("解析代理URL失败: %v", err)
		}

		switch parsedURL.Scheme {
		case "http", "https":
			// HTTP/HTTPS 代理
			transport.Proxy = http.ProxyURL(parsedURL)
		case "socks5":
			// SOCKS5 代理
			auth := &proxy.Auth{}
			if parsedURL.User != nil {
				auth.User = parsedURL.User.Username()
				auth.Password, _ = parsedURL.User.Password()
			}

			dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
			if err != nil {
				return nil, fmt.Errorf("创建SOCKS5代理失败: %v", err)
			}

			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		default:
			return nil, fmt.Errorf("不支持的代理协议: %s (支持: http, https, socks5)", parsedURL.Scheme)
		}
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	return client, nil
}

// Create 创建 RT
func (s *rtService) Create(rt *model.RT) error {
	// 如果name为空，生成32位UUID
	if strings.TrimSpace(rt.BizId) == "" {
		rt.BizId = generateRandomID()
	}

	// 检查name是否已存在
	existing, err := s.repo.GetByBizId(rt.BizId)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("RT名称 '%s' 已存在", rt.BizId)
	}

	// 检查token是否已存在
	existingToken, err := s.repo.GetByToken(rt.Rt)
	if err != nil {
		return err
	}
	if existingToken != nil {
		return fmt.Errorf("此RT Token已存在")
	}

	// 填充默认配置值（如果字段为空）
	cfg := config.Get()
	if rt.ClientID == "" {
		rt.ClientID = cfg.OpenAI.ClientID
		logger.Info("创建RT时填充默认 client_id", "client_id", rt.ClientID)
	}
	if rt.Proxy == "" && cfg.OpenAI.Proxy != "" {
		rt.Proxy = cfg.OpenAI.Proxy
		logger.Info("创建RT时填充默认 proxy", "proxy", rt.Proxy)
	}

	return s.repo.Create(rt)
}

// Update 更新 RT
func (s *rtService) Update(id int64, updates map[string]interface{}) (*model.RT, error) {
	rt, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, fmt.Errorf("RT不存在")
	}

	// 应用更新
	if bizId, ok := updates["biz_id"].(string); ok && bizId != "" {
		// 检查新名称是否与其他RT冲突
		if bizId != rt.BizId {
			existing, _ := s.repo.GetByBizId(bizId)
			if existing != nil {
				return nil, fmt.Errorf("RT名称 '%s' 已被使用", bizId)
			}
		}
		rt.BizId = bizId
	}
	if proxy, ok := updates["proxy"].(string); ok {
		rt.Proxy = proxy
	}
	if tag, ok := updates["tag"].(string); ok {
		rt.Tag = tag
	}
	if enabled, ok := updates["enabled"].(bool); ok {
		rt.Enabled = enabled
	}
	if memo, ok := updates["memo"].(string); ok {
		rt.Memo = memo
	}

	if err := s.repo.Update(rt); err != nil {
		return nil, err
	}

	return rt, nil
}

// GetByID 获取RT
func (s *rtService) GetByID(id int64) (*model.RT, error) {
	return s.repo.GetByID(id)
}

// GetByBizId 根据业务ID获取RT
func (s *rtService) GetByBizId(bizId string) (*model.RT, error) {
	return s.repo.GetByBizId(bizId)
}

// GetByEmail 根据邮箱获取RT
func (s *rtService) GetByEmail(email string) (*model.RT, error) {
	return s.repo.GetByEmail(email)
}

// List 获取列表
func (s *rtService) List(page, pageSize int, name string, tag string, email string, typeStr string, enabled *bool, createDate string) ([]*model.RT, int64, error) {
	return s.repo.List(page, pageSize, name, tag, email, typeStr, enabled, createDate)
}

// Delete 删除RT
func (s *rtService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// BatchDelete 批量删除
func (s *rtService) BatchDelete(ids []int64) (int, int, error) {
	return s.repo.BatchDelete(ids)
}

// OpenAI Token 响应结构
type OpenAITokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// OpenAI Error 响应结构
type OpenAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Account Check 响应结构
type AccountCheckResponse struct {
	Accounts map[string]struct {
		Account struct {
			PlanType string `json:"plan_type"`
		} `json:"account"`
	} `json:"accounts"`
	AccountOrdering []string `json:"account_ordering"`
}

// User Info 响应结构
type UserInfoResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Refresh 刷新单个RT
func (s *rtService) Refresh(id int64, refreshUserInfo, refreshAccountInfo bool) (*model.RT, error) {
	rt, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, fmt.Errorf("RT不存在")
	}

	logger.Info("开始刷新RT", "id", id, "name", rt.BizId, "has_proxy", rt.Proxy != "")

	// 获取 client_id，优先使用 RT 记录中的，如果为空则使用配置文件中的默认值
	clientID := rt.ClientID
	if clientID == "" {
		cfg := config.Get()
		clientID = cfg.OpenAI.ClientID
		logger.Info("使用配置文件中的默认 client_id", "client_id", clientID)
	}

	// 构造请求体
	requestBody := map[string]string{
		"client_id":     clientID,
		"grant_type":    "refresh_token",
		"redirect_uri":  "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
		"refresh_token": rt.Rt,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logger.Error("构造请求体失败", "error", err)
		return nil, fmt.Errorf("构造请求体失败: %v", err)
	}

	// 创建支持 SOCKS5 的 HTTP 客户端
	client, err := createHTTPClient(rt.Proxy, 10*time.Second)
	if err != nil {
		logger.Warn("创建HTTP客户端失败", "proxy", rt.Proxy, "error", err)
		// 使用无代理的客户端
		client = &http.Client{Timeout: 10 * time.Second}
	} else if rt.Proxy != "" {
		logger.Info("使用代理", "proxy", rt.Proxy)
	}

	// 创建请求
	req, err := http.NewRequest("POST", "https://auth.openai.com/oauth/token", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("创建请求失败", "error", err)
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("请求失败", "error", err)
		// 保存失败结果
		rt.RefreshResult = fmt.Sprintf("请求失败: %v", err)
		s.repo.Update(rt)
		return rt, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取响应失败", "error", err)
		rt.RefreshResult = fmt.Sprintf("读取响应失败: %v", err)
		s.repo.Update(rt)
		return rt, fmt.Errorf("读取响应失败: %v", err)
	}

	// 保存完整的响应body到refresh_result
	rt.RefreshResult = string(body)

	// 解析响应
	if resp.StatusCode == 200 {
		// 成功响应
		var tokenResp OpenAITokenResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			logger.Error("解析成功响应失败", "error", err, "body", string(body))
			s.repo.Update(rt)
			return rt, fmt.Errorf("解析响应失败: %v", err)
		}

		// 保存旧的RT Token到LastRT
		rt.LastRT = rt.Rt
		// 更新为新的RT Token
		rt.Rt = tokenResp.RefreshToken
		// 保存Access Token
		rt.At = tokenResp.AccessToken
		// 更新刷新时间
		now := time.Now()
		rt.LastRefreshTime = &now

		logger.Info("刷新RT成功",
			"id", id,
			"name", rt.BizId,
			"old_rt", rt.LastRT[:20]+"...",
			"new_rt", rt.Rt[:20]+"...",
		)

		// 根据参数决定是否获取用户信息
		if refreshUserInfo {
			logger.Info("UserInfo为空，开始获取用户信息", "id", id, "name", rt.BizId)
			if err := s.fetchUserInfo(rt); err != nil {
				logger.Warn("获取用户信息失败", "id", id, "name", rt.BizId, "error", err)
				// 不影响刷新流程，继续执行
			}
		}

		// 根据参数决定是否获取账号信息
		if refreshAccountInfo {
			logger.Info("AccountInfo为空，开始获取账号信息", "id", id, "name", rt.BizId)
			if err := s.fetchAccountInfo(rt); err != nil {
				logger.Warn("获取账号信息失败", "id", id, "name", rt.BizId, "error", err)
				// 不影响刷新流程，继续执行
			}
		}
	} else {
		// 失败响应
		var errorResp OpenAIErrorResponse
		var errorMsg string

		if err := json.Unmarshal(body, &errorResp); err != nil {
			// 解析失败，使用原始响应
			logger.Error("解析错误响应失败", "error", err, "status", resp.StatusCode, "body", string(body))
			errorMsg = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))
		} else {
			logger.Error("刷新RT失败",
				"id", id,
				"name", rt.BizId,
				"error_code", errorResp.Error.Code,
				"error_message", errorResp.Error.Message,
			)
			errorMsg = fmt.Sprintf("%s: %s", errorResp.Error.Code, errorResp.Error.Message)
		}

		// RT、LastRT和Enabled保持不变，只更新RefreshResult
		// 更新数据库
		if err := s.repo.Update(rt); err != nil {
			logger.Error("更新RT失败", "error", err)
		}

		// 返回错误给调用方
		return rt, fmt.Errorf("刷新失败: %s", errorMsg)
	}

	// 成功时才更新数据库
	if err := s.repo.Update(rt); err != nil {
		logger.Error("更新RT失败", "error", err)
		return nil, fmt.Errorf("更新RT失败: %v", err)
	}

	return rt, nil
}

// createTLSClient 创建带有 Firefox TLS 指纹的客户端
func createTLSClient(proxyURL string, timeout time.Duration) (tls_client.HttpClient, error) {
	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Firefox_133), // 使用 Firefox 133 指纹
		tls_client.WithTimeoutSeconds(int(timeout.Seconds())),
	}

	// 如果有代理，添加代理配置
	if proxyURL != "" {
		options = append(options, tls_client.WithProxyUrl(proxyURL))
	}

	return tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
}

// fetchUserInfo 获取用户信息
func (s *rtService) fetchUserInfo(rt *model.RT) error {
	// 创建 TLS 客户端（使用 Firefox 指纹）
	client, err := createTLSClient(rt.Proxy, 30*time.Second)
	if err != nil {
		logger.Warn("创建TLS客户端失败", "proxy", rt.Proxy, "error", err)
		// 尝试不使用代理
		client, err = createTLSClient("", 30*time.Second)
		if err != nil {
			return fmt.Errorf("创建TLS客户端失败: %v", err)
		}
	}

	// 创建请求（使用 fhttp）
	req, err := http2.NewRequest("GET", "https://chatgpt.com/backend-api/me", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头（模拟 Firefox 浏览器行为）
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("authorization", "Bearer "+rt.At)
	req.Header.Set("dnt", "1")
	req.Header.Set("oai-language", "zh-CN")
	req.Header.Set("priority", "u=1")
	req.Header.Set("referer", "https://chatgpt.com/")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("【刷新用户信息】请求失败", "id", rt.ID, "error", err)
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	logger.Info(string(body))

	// 解析响应
	if resp.StatusCode == 200 {
		// 保存完整的用户信息到 UserInfo
		rt.UserInfo = string(body)

		var userResp UserInfoResponse
		if err := json.Unmarshal(body, &userResp); err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}

		// 保存Email和UserName
		if userResp.Email != "" {
			rt.Email = userResp.Email
		}
		if userResp.Name != "" {
			rt.UserName = userResp.Name
		}

		if rt.Email != "" || rt.UserName != "" {
			logger.Info("获取用户信息成功", "id", rt.ID, "biz_id", rt.BizId, "user_name", rt.UserName, "email", rt.Email)
		} else {
			logger.Warn("用户邮箱和名称均为空", "id", rt.ID, "biz_id", rt.BizId)
		}
	} else {
		logger.Warn("获取用户信息失败", "id", rt.ID, "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// fetchAccountInfo 获取账号信息
func (s *rtService) fetchAccountInfo(rt *model.RT) error {
	// 创建 TLS 客户端（使用 Firefox 指纹）
	client, err := createTLSClient(rt.Proxy, 30*time.Second)
	if err != nil {
		logger.Warn("创建TLS客户端失败", "proxy", rt.Proxy, "error", err)
		// 尝试不使用代理
		client, err = createTLSClient("", 30*time.Second)
		if err != nil {
			return fmt.Errorf("创建TLS客户端失败: %v", err)
		}
	}

	// 创建请求（使用 fhttp）
	req, err := http2.NewRequest("GET", "https://chatgpt.com/backend-api/accounts/check/v4-2023-04-27?timezone_offset_min=-480", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头（模拟 Firefox 浏览器行为）
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("authorization", "Bearer "+rt.At)
	req.Header.Set("dnt", "1")
	req.Header.Set("oai-language", "zh-CN")
	req.Header.Set("priority", "u=1")
	req.Header.Set("referer", "https://chatgpt.com/")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("【刷新账号信息】请求失败", "id", rt.ID, "error", err)
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	if resp.StatusCode == 200 {
		// 保存完整的账号信息到 AccountInfo
		rt.AccountInfo = string(body)

		var accountResp AccountCheckResponse
		if err := json.Unmarshal(body, &accountResp); err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}

		// 遍历所有账号，寻找最佳的plan_type
		var selectedPlanType string
		var hasFree bool

		for accountID, accountData := range accountResp.Accounts {
			planType := accountData.Account.PlanType
			logger.Info("发现账号", "account_id", accountID, "plan_type", planType)

			if planType == "free" {
				hasFree = true
			} else if planType != "" {
				// 找到非free的plan_type，优先使用
				selectedPlanType = planType
				logger.Info("选择plan_type", "account_id", accountID, "plan_type", planType)
				break
			}
		}

		// 如果没有找到非free的，但有free，则使用free
		if selectedPlanType == "" && hasFree {
			selectedPlanType = "free"
		}

		// 保存到Type字段
		if selectedPlanType != "" {
			rt.Type = selectedPlanType
			logger.Info("设置账号类型", "id", rt.ID, "name", rt.BizId, "type", rt.Type)
		} else {
			logger.Warn("未找到任何plan_type", "id", rt.ID, "name", rt.BizId)
		}
	} else {
		logger.Warn("获取账号信息失败", "id", rt.ID, "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// RefreshUserInfo 刷新用户信息
func (s *rtService) RefreshUserInfo(id int64) (*model.RT, error) {
	rt, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, fmt.Errorf("RT不存在")
	}

	// 检查是否有 AT
	if rt.At == "" {
		return nil, fmt.Errorf("AT为空，无法获取用户信息")
	}

	logger.Info("开始刷新用户信息", "id", id, "biz_id", rt.BizId)

	// 获取用户信息
	if err := s.fetchUserInfo(rt); err != nil {
		logger.Error("刷新用户信息失败", "id", id, "biz_id", rt.BizId, "error", err)
		return nil, fmt.Errorf("刷新用户信息失败: %v", err)
	}

	// 更新数据库
	if err := s.repo.Update(rt); err != nil {
		logger.Error("保存用户信息失败", "id", id, "error", err)
		return nil, fmt.Errorf("保存用户信息失败: %v", err)
	}

	logger.Info("刷新用户信息成功", "id", id, "biz_id", rt.BizId, "user_name", rt.UserName, "email", rt.Email)
	return rt, nil
}

// RefreshAccountInfo 刷新账号信息
func (s *rtService) RefreshAccountInfo(id int64) (*model.RT, error) {
	rt, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, fmt.Errorf("RT不存在")
	}

	// 检查是否有 AT
	if rt.At == "" {
		return nil, fmt.Errorf("AT为空，无法获取账号信息")
	}

	logger.Info("开始刷新账号信息", "id", id, "biz_id", rt.BizId)

	// 获取账号信息
	if err := s.fetchAccountInfo(rt); err != nil {
		logger.Error("刷新账号信息失败", "id", id, "biz_id", rt.BizId, "error", err)
		return nil, fmt.Errorf("刷新账号信息失败: %v", err)
	}

	// 更新数据库
	if err := s.repo.Update(rt); err != nil {
		logger.Error("保存账号信息失败", "id", id, "error", err)
		return nil, fmt.Errorf("保存账号信息失败: %v", err)
	}

	logger.Info("刷新账号信息成功", "id", id, "biz_id", rt.BizId, "type", rt.Type)
	return rt, nil
}

// BatchRefresh 批量刷新
func (s *rtService) BatchRefresh(ids []int64) (int, int, []map[string]interface{}, error) {
	rts, err := s.repo.GetByIDs(ids)
	if err != nil {
		return 0, 0, nil, err
	}

	successCount := 0
	failCount := 0
	results := make([]map[string]interface{}, 0, len(rts))

	for i, rt := range rts {
		result := map[string]interface{}{
			"rt_name": rt.BizId,
		}

		// 刷新，批量刷新时默认获取用户信息和账号信息
		_, err := s.Refresh(rt.ID, true, true)
		if err != nil {
			failCount++
			result["success"] = false
			result["message"] = err.Error()
		} else {
			successCount++
			result["success"] = true
			result["message"] = "刷新成功"
		}

		results = append(results, result)

		// 相邻刷新之间随机延迟 1-3 秒（最后一个不需要延迟）
		if i < len(rts)-1 {
			delay := time.Duration(1+rand.Intn(3)) * time.Second
			logger.Info("批量刷新延迟", "delay_seconds", delay.Seconds())
			time.Sleep(delay)
		}
	}

	return successCount, failCount, results, nil
}

// BatchImport 批量导入（batchName参数已弃用，每个RT都会生成唯一的32位UUID）
func (s *rtService) BatchImport(batchName string, tag string, tokens []string, proxyList []string) (int, int, error) {
	successCount := 0
	failCount := 0

	// 去重
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token != "" {
			uniqueTokens[token] = true
		}
	}

	for token := range uniqueTokens {
		// 检查token是否已存在
		existing, _ := s.repo.GetByToken(token)
		if existing != nil {
			failCount++
			logger.Warn("Token已存在，跳过", "token", token[:20]+"...")
			continue
		}

		// 生成唯一的32位UUID
		name := generateRandomID()
		existingName, _ := s.repo.GetByBizId(name)
		for existingName != nil {
			name = generateRandomID()
			existingName, _ = s.repo.GetByBizId(name)
		}

		// 随机选择代理
		var proxy string
		if len(proxyList) > 0 {
			proxy = proxyList[rand.Intn(len(proxyList))]
		}

		// 创建RT
		rt := &model.RT{
			BizId:   name,
			Rt:      token,
			Proxy:   proxy,
			Tag:     tag,
			Enabled: false,
		}

		// 填充默认配置值（如果字段为空）
		cfg := config.Get()
		if rt.ClientID == "" {
			rt.ClientID = cfg.OpenAI.ClientID
		}
		if rt.Proxy == "" && cfg.OpenAI.Proxy != "" {
			rt.Proxy = cfg.OpenAI.Proxy
		}

		if err := s.repo.Create(rt); err != nil {
			failCount++
			logger.Error("创建RT失败", "name", name, "error", err)
		} else {
			successCount++
			logger.Info("导入RT成功", "name", name)
		}
	}

	return successCount, failCount, nil
}

// getProxyListFromConfig 从配置中获取代理列表
func (s *rtService) getProxyListFromConfig() []string {
	config, err := s.configRepo.GetByKey("proxy_list")
	if err != nil || config == nil {
		return nil
	}

	var proxyList []string
	if err := json.Unmarshal([]byte(config.ConfigValue), &proxyList); err != nil {
		return nil
	}

	return proxyList
}

// AutoRefreshAll 自动刷新所有启用的RT
func (s *rtService) AutoRefreshAll() error {
	// 获取所有启用的RT
	enabled := true
	rts, count, err := s.List(1, 10000, "", "", "", "", &enabled, "")
	if err != nil {
		return err
	}

	logger.Info("开始自动刷新", "count", count)

	// 提取RT的ID
	ids := make([]int64, 0, len(rts))
	for _, rt := range rts {
		ids = append(ids, rt.ID)
	}

	// 批量刷新
	if len(ids) > 0 {
		successCount, failCount, _, err := s.BatchRefresh(ids)
		if err != nil {
			logger.Error("批量刷新失败", "error", err)
			return err
		}
		logger.Info("自动刷新完成", "success", successCount, "fail", failCount)
	}

	return nil
}

// generateRandomID 生成32位UUID（去掉破折号）
func generateRandomID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
