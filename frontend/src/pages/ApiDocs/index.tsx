import React, { useState } from 'react';
import { Card, Collapse, Typography, Tag, Space, Input, Button, message } from 'antd';
import { ApiOutlined, CopyOutlined, CheckOutlined } from '@ant-design/icons';

const { Panel } = Collapse;
const { Title, Paragraph, Text } = Typography;
const { TextArea } = Input;

const ApiDocs: React.FC = () => {
  const [copiedIndex, setCopiedIndex] = useState<string | null>(null);

  // 复制到剪贴板
  const copyToClipboard = (text: string, index: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopiedIndex(index);
      message.success('已复制到剪贴板');
      setTimeout(() => setCopiedIndex(null), 2000);
    });
  };

  // 对外公开API接口列表（仅3个核心接口）
  const apiList = [
    {
      category: '对外公开 API',
      apis: [
        {
          name: '刷新RT并获取AT',
          method: 'POST',
          path: '/public-api/refresh',
          description: '根据 biz_id 或 email 查找RT并刷新，返回新的 access_token。优先级：biz_id > email',
          needAuth: true,
          curl: `# 使用 biz_id 查询
curl -X POST http://localhost:8080/public-api/refresh \\
  -H "Content-Type: application/json" \\
  -H "X-API-Secret: my-api-secret-2025" \\
  -d '{
    "biz_id": "user001"
  }'

# 使用 email 查询
curl -X POST http://localhost:8080/public-api/refresh \\
  -H "Content-Type: application/json" \\
  -H "X-API-Secret: my-api-secret-2025" \\
  -d '{
    "email": "user@example.com"
  }'

# 响应示例
{
  "success": true,
  "msg": "刷新成功",
  "data": {
    "biz_id": "user001",
    "email": "user@example.com",
    "access_token": "eyJhbGciOiJSUzI1NiIsImt...",
    "refresh_token": "rt_xxx...",
    "type": "team",
    "user_name": "John Doe"
  }
}`,
        },
        {
          name: '获取AT（不刷新）',
          method: 'POST',
          path: '/public-api/get-at',
          description: '根据 biz_id 或 email 查找RT，直接返回已有的 access_token，不执行刷新操作',
          needAuth: true,
          curl: `# 使用 biz_id 查询
curl -X POST http://localhost:8080/public-api/get-at \\
  -H "Content-Type: application/json" \\
  -H "X-API-Secret: my-api-secret-2025" \\
  -d '{
    "biz_id": "user001"
  }'

# 使用 email 查询
curl -X POST http://localhost:8080/public-api/get-at \\
  -H "Content-Type: application/json" \\
  -H "X-API-Secret: my-api-secret-2025" \\
  -d '{
    "email": "user@example.com"
  }'

# 响应示例
{
  "success": true,
  "msg": "获取成功",
  "data": {
    "biz_id": "user001",
    "email": "user@example.com",
    "access_token": "eyJhbGciOiJSUzI1NiIsImt...",
    "refresh_token": "rt_xxx...",
    "type": "team",
    "user_name": "John Doe"
  }
}`,
        },
        {
          name: '健康检查',
          method: 'GET',
          path: '/public-api/health',
          description: '检查服务健康状态，无需认证',
          needAuth: true,
          curl: `curl -X GET http://localhost:8080/public-api/health \\
  -H "X-API-Secret: my-api-secret-2025"

# 响应示例
{
  "success": true,
  "msg": "服务运行正常"
}`,
        },
      ],
    },
  ];

  return (
    <div style={{ padding: '0px' }}>
      <Card
        title={
          <Space>
            <ApiOutlined style={{ fontSize: 20, color: '#1890ff' }} />
            <span style={{ fontSize: 18, fontWeight: 'bold' }}>API 接口文档</span>
          </Space>
        }
        extra={
          <Tag color="green">对外公开 API（使用 API Secret 认证）</Tag>
        }
      >
        <Paragraph type="secondary">
          <Text strong>认证方式：</Text>
          <ul>
            <li>
              <Text strong>API Secret 认证：</Text> 使用配置文件中的固定密钥进行认证
              <br />
              <Text code>X-API-Secret: my-api-secret-2025</Text>
              <br />
              <Text type="warning">（api_secret 可在 config.yaml 的 auth.api_secret 中配置，默认值：my-api-secret-2025）</Text>
            </li>
          </ul>
        </Paragraph>

        <Paragraph type="secondary">
          <Text strong>API 规范说明：</Text>
          <ul>
            <li>所有接口使用 <Text code>POST</Text> 方法（健康检查除外，使用GET）</li>
            <li>请求参数放在 <Text code>JSON Body</Text> 中</li>
            <li>响应格式统一为：<Text code>{`{ "success": boolean, "msg": string, "data": any }`}</Text></li>
            <li>认证接口需要在 Header 中添加：<Text code>X-API-Secret: YOUR_API_SECRET</Text></li>
          </ul>
        </Paragraph>

        <Paragraph type="secondary">
          <Text strong>查询优先级：</Text>
          <ul>
            <li><Text code>biz_id</Text> 优先级最高：如果提供了 biz_id，则使用 biz_id 查询</li>
            <li><Text code>email</Text> 优先级次之：如果未提供 biz_id，则使用 email 查询</li>
            <li>两者必须至少提供一个，否则返回参数错误</li>
          </ul>
        </Paragraph>

        {apiList.map((category, catIndex) => (
          <Card
            key={catIndex}
            type="inner"
            title={<Text strong>{category.category}</Text>}
            style={{ marginBottom: 16 }}
          >
            <Collapse accordion>
              {category.apis.map((api, apiIndex) => {
                const index = `${catIndex}-${apiIndex}`;
                return (
                  <Panel
                    header={
                      <Space>
                        <Tag color={api.method === 'POST' ? 'green' : api.method === 'GET' ? 'blue' : 'orange'}>
                          {api.method}
                        </Tag>
                        <Text strong>{api.name}</Text>
                        <Text type="secondary">{api.path}</Text>
                        {api.needAuth && <Tag color="gold">需要认证</Tag>}
                      </Space>
                    }
                    key={index}
                  >
                    <Paragraph>{api.description}</Paragraph>
                    <div style={{ position: 'relative' }}>
                      <Button
                        type="primary"
                        size="small"
                        icon={copiedIndex === index ? <CheckOutlined /> : <CopyOutlined />}
                        onClick={() => copyToClipboard(api.curl, index)}
                        style={{
                          position: 'absolute',
                          right: 8,
                          top: 8,
                          zIndex: 1,
                        }}
                      >
                        {copiedIndex === index ? '已复制' : '复制'}
                      </Button>
                      <TextArea
                        value={api.curl}
                        readOnly
                        autoSize={{ minRows: 3, maxRows: 15 }}
                        style={{
                          fontFamily: 'monospace',
                          fontSize: 12,
                          backgroundColor: '#f5f5f5',
                        }}
                      />
                    </div>
                  </Panel>
                );
              })}
            </Collapse>
          </Card>
        ))}
      </Card>
    </div>
  );
};

export default ApiDocs;
