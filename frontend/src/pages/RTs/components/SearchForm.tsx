import React from 'react';
import { Card, Form, Input, Select, DatePicker, Button, Row, Col } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import type { FormInstance } from 'antd/es/form';

interface SearchFormProps {
  form: FormInstance;
  onSearch: () => void;
  onReset: () => void;
}

const SearchForm: React.FC<SearchFormProps> = ({
  form,
  onSearch,
  onReset,
}) => {
  return (
    <Card 
      size="small" 
      style={{ marginBottom: 16, background: '#fafafa' }}
    >
      <Form form={form} layout="inline" onFinish={onSearch}>
        <Row gutter={16} style={{ width: '100%' }}>
          <Col xs={24} sm={12} md={6}>
            <Form.Item name="biz_id" label="业务ID" style={{ marginBottom: 8 }}>
              <Input placeholder="输入业务ID" allowClear />
            </Form.Item>
          </Col>
          
          <Col xs={24} sm={12} md={6}>
            <Form.Item name="tag" label="标签" style={{ marginBottom: 8 }}>
              <Input placeholder="输入标签" allowClear />
            </Form.Item>
          </Col>

          <Col xs={24} sm={12} md={6}>
            <Form.Item name="email" label="邮箱" style={{ marginBottom: 8 }}>
              <Input placeholder="输入邮箱" allowClear />
            </Form.Item>
          </Col>

          <Col xs={24} sm={12} md={6}>
            <Form.Item name="type" label="类型" style={{ marginBottom: 8 }}>
              <Input placeholder="输入类型" allowClear />
            </Form.Item>
          </Col>

          <Col xs={24} sm={12} md={6}>
            <Form.Item name="enabled" label="状态" style={{ marginBottom: 8 }}>
              <Select placeholder="选择状态" allowClear>
                <Select.Option value={true}>启用</Select.Option>
                <Select.Option value={false}>禁用</Select.Option>
              </Select>
            </Form.Item>
          </Col>

          <Col xs={24} sm={12} md={6}>
            <Form.Item name="createDate" label="创建日期" style={{ marginBottom: 8 }}>
              <DatePicker 
                placeholder="选择创建日期" 
                style={{ width: '100%' }}
                allowClear
              />
            </Form.Item>
          </Col>
        </Row>

        <Row style={{ marginTop: 8 }}>
          <Col span={24}>
            <Button type="primary" icon={<SearchOutlined />} onClick={onSearch}>
              搜索
            </Button>
            <Button onClick={onReset} style={{ marginLeft: 8 }}>
              重置
            </Button>
          </Col>
        </Row>
      </Form>
    </Card>
  );
};

export default SearchForm;

