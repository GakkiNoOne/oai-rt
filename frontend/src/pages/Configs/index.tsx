import React, { useState, useEffect } from 'react';
import {
  Card,
  Button,
  Space,
  Form,
  Input,
  InputNumber,
  Typography,
  message,
  Alert,
  Spin,
  Switch,
  Modal,
  Row,
  Col
} from 'antd';
import {
  SaveOutlined,
  SettingOutlined
} from '@ant-design/icons';
import { configsApi, SystemConfigsResponse } from '@/api/configs';

const { Title, Text } = Typography;
const { TextArea } = Input;

const ConfigManagement: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  // 加载配置
  const loadConfigs = async () => {
    setLoading(true);
    try {
      const response = await configsApi.getSystemConfigs();

      if (response.success && response.data) {
        const { configs } = response.data;
        
        // 解析 JSON 字段
        const proxyList = configs.proxy_list ? JSON.parse(configs.proxy_list) : [];
        
        // 解析 clientId 列表
        const clientIdList = configs.client_id_list ? JSON.parse(configs.client_id_list) : [];
        
        // 如果 clientId 列表为空，使用默认值
        const clientIdValue = clientIdList.length > 0 
          ? clientIdList.join('\n') 
          : 'app_WXrF1LSkiTtfYqiL6XtjygvX';
        
        form.setFieldsValue({
          proxy_list: proxyList.join('\n'),
          client_id_list: clientIdValue,
          auto_refresh_enabled: configs.auto_refresh_enabled === 'true',
          auto_refresh_interval: parseInt(configs.auto_refresh_interval) || 2,
        });
      }
    } catch (error) {
      console.error('加载配置失败:', error);
      message.error('加载配置失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadConfigs();
  }, []);

  // 保存配置
  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      
      // 显示确认对话框
      Modal.confirm({
        title: '确认保存配置',
        content: (
          <div>
            <p>保存配置后，系统将自动执行以下操作：</p>
            <ul style={{ marginTop: 8, paddingLeft: 20 }}>
              <li>如果<strong>代理列表</strong>或<strong>Client ID 列表</strong>发生变化，将自动更新所有 RT 的配置</li>
              <li>不在新列表中的 RT，将随机分配新的代理和 Client ID</li>
              <li>自动刷新配置将在下次调度时生效</li>
            </ul>
            <Alert 
              message="提示" 
              description="代理可以设置为空（所有 RT 将使用直连），Client ID 必须至少保留一个"
              type="info" 
              showIcon 
              style={{ marginTop: 12 }}
            />
          </div>
        ),
        okText: '确定保存',
        cancelText: '取消',
        width: 560,
        onOk: async () => {
          setSaving(true);
          try {
            // 处理代理列表
            const proxyList = values.proxy_list
              ? values.proxy_list
                  .split('\n')
                  .map((line: string) => line.trim())
                  .filter((line: string) => line !== '')
              : [];

            // 处理 clientId 列表
            const clientIdList = values.client_id_list
              ? values.client_id_list
                  .split('\n')
                  .map((line: string) => line.trim())
                  .filter((line: string) => line !== '')
              : [];

            // 构建配置对象
            const configs: Record<string, string> = {
              proxy_list: JSON.stringify(proxyList),
              client_id_list: JSON.stringify(clientIdList),
              auto_refresh_enabled: values.auto_refresh_enabled ? 'true' : 'false',
              auto_refresh_interval: values.auto_refresh_interval.toString(),
            };

            const response = await configsApi.saveSystemConfigs(configs);

            if (response.success) {
              message.success('保存成功，已自动更新相关 RT 配置');
            }
          } catch (error) {
            console.error('保存失败:', error);
            message.error('保存失败');
          } finally {
            setSaving(false);
          }
        },
      });
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  return (
    <div>
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: 24
      }}>
        <Space align="center">
          <SettingOutlined style={{ fontSize: 24, color: '#1890ff' }} />
          <Title level={2} style={{ margin: 0 }}>系统配置</Title>
        </Space>
      </div>

      <Spin spinning={loading}>
        <Form
          form={form}
          layout="vertical"
        >
          <Row gutter={24}>
            {/* 左侧列 */}
            <Col xs={24} lg={12}>
              {/* 代理配置 */}
              <Card 
                title={<Text strong style={{ fontSize: 16 }}>🔌 代理服务器配置</Text>}
                size="small"
                style={{ marginBottom: 24 }}
              >
                <Alert
                  message="配置说明"
                  description={
                    <div>
                      <div>• 代理<Text strong>可选</Text>，不配置则使用本机 IP 发送请求</div>
                      <div>• 请求时会从列表中<Text strong>随机选择</Text>一个代理</div>
                      <div>• <Text strong type="danger">每行一个</Text>，支持配置<Text strong>多个</Text>代理</div>
                      <div>• 支持协议：<Text code>http://</Text>、<Text code>https://</Text>、<Text code>socks5://</Text></div>
                    </div>
                  }
                  type="info"
                  showIcon
                  style={{ marginBottom: 16 }}
                />

                <Form.Item
                  name="proxy_list"
                  label={<Text strong>代理列表（每行一个，可选）</Text>}
                  rules={[
                    {
                      validator: (_, value) => {
                        if (!value || value.trim() === '') {
                          return Promise.resolve();
                        }
                        const lines = value.split('\n').filter((line: string) => line.trim() !== '');
                        if (lines.length === 0) {
                          return Promise.resolve();
                        }
                        for (const line of lines) {
                          if (!line.match(/^(https?|socks5):\/\/.+/)) {
                            return Promise.reject(new Error(`代理地址格式错误: ${line}`));
                          }
                        }
                        return Promise.resolve();
                      }
                    }
                  ]}
                  style={{ marginBottom: 0 }}
                >
                  <TextArea
                    rows={6}
                    placeholder="每行一个代理地址（可选），格式示例：&#10;http://127.0.0.1:7890&#10;https://proxy.example.com:8080&#10;socks5://user:pass@host:port"
                    style={{ fontFamily: 'monospace', fontSize: '12px' }}
                  />
                </Form.Item>
              </Card>

              {/* 自动刷新配置 */}
              <Card 
                title={<Text strong style={{ fontSize: 16 }}>🔄 自动刷新配置</Text>}
                size="small"
                style={{ marginBottom: 24 }}
              >
                <Form.Item
                  name="auto_refresh_enabled"
                  label={<Text strong>启用自动刷新</Text>}
                  valuePropName="checked"
                  initialValue={false}
                  style={{ marginBottom: 16 }}
                >
                  <Switch 
                    checkedChildren="开启" 
                    unCheckedChildren="关闭"
                  />
                </Form.Item>
                <Alert
                  message="自动刷新功能会定期检查所有启用的 RT 的有效性，失效的 RT 将被自动禁用。"
                  type="info"
                  showIcon
                  style={{ marginBottom: 16 }}
                />

                <Form.Item
                  name="auto_refresh_interval"
                  label={<Text strong>刷新间隔（天）</Text>}
                  rules={[
                    { required: true, message: '请输入刷新间隔' },
                    { type: 'number', min: 1, max: 30, message: '刷新间隔必须在 1-30 天之间' }
                  ]}
                  extra="建议设置 2 天，避免频繁刷新"
                  style={{ marginBottom: 0 }}
                >
                  <InputNumber
                    style={{ width: '100%' }}
                    min={1}
                    max={30}
                    placeholder="请输入刷新间隔（天）"
                  />
                </Form.Item>
              </Card>
            </Col>

            {/* 右侧列 */}
            <Col xs={24} lg={12}>
              {/* Client ID 配置 */}
              <Card 
                title={<Text strong style={{ fontSize: 16 }}>🔑 Client ID 配置</Text>}
                size="small"
                style={{ marginBottom: 24 }}
              >
                <Alert
                  message="配置说明"
                  description={
                    <div>
                      <div>• Client ID 用于刷新 RT</div>
                      <div>• 系统会从列表中<Text strong>随机选择</Text>一个 Client ID</div>
                      <div>• <Text strong type="danger">每行一个</Text>，支持配置<Text strong>多个</Text> Client ID</div>
                      <div>• 如果不配置，使用默认 Client ID</div>
                    </div>
                  }
                  type="info"
                  showIcon
                  style={{ marginBottom: 16 }}
                />

                <Form.Item
                  name="client_id_list"
                  label={<Text strong>Client ID 列表（每行一个，可选）</Text>}
                  style={{ marginBottom: 0 }}
                >
                  <TextArea
                    rows={6}
                    placeholder="每行一个 Client ID（可选），格式示例：&#10;app_WXrF1LSkiTtfYqiL6XtjygvX&#10;app_AnotherClientId123456&#10;如不配置，使用系统默认值"
                    style={{ fontFamily: 'monospace', fontSize: '12px' }}
                  />
                </Form.Item>
              </Card>
            </Col>
          </Row>

          {/* 底部操作按钮 */}
          <Form.Item style={{ marginTop: 32 }}>
            <Space>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                onClick={handleSave}
                loading={saving}
                size="large"
              >
                保存配置
              </Button>
              <Button
                size="large"
                onClick={loadConfigs}
              >
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Spin>
    </div>
  );
};

export default ConfigManagement;
