import React from 'react';
import { Modal, Form, Input, Select, AutoComplete } from 'antd';
import type { FormInstance } from 'antd/es/form';

interface BatchImportModalProps {
  open: boolean;
  form: FormInstance;
  loading: boolean;
  proxyList: string[];
  loadingProxy: boolean;
  clientIDList: string[];
  loadingClientID: boolean;
  onOk: () => void;
  onCancel: () => void;
}

const BatchImportModal: React.FC<BatchImportModalProps> = ({
  open,
  form,
  loading,
  proxyList,
  loadingProxy,
  clientIDList,
  loadingClientID,
  onOk,
  onCancel,
}) => {
  return (
    <Modal
      title="批量导入RT"
      open={open}
      onOk={onOk}
      onCancel={onCancel}
      width={800}
      okText="导入"
      cancelText="取消"
      confirmLoading={loading}
    >
      <Form
        form={form}
        layout="vertical"
      >
        <Form.Item
          name="tag"
          label="标签"
          extra="可选，为这批RT设置统一标签"
        >
          <Input placeholder="例如：测试账号、生产环境等（可选）" allowClear />
        </Form.Item>

        <Form.Item
          name="proxy"
          label="代理"
          extra="可选，为这批RT设置统一代理，不配置则使用本机IP"
        >
          <Select
            placeholder="从列表选择代理（可选）"
            loading={loadingProxy}
            showSearch
            optionFilterProp="children"
            allowClear
          >
            {proxyList.map((proxy, index) => (
              <Select.Option key={index} value={proxy}>
                {proxy}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="client_id"
          label="Client ID"
          extra="可选，为这批RT设置统一的Client ID，不配置则使用默认值"
        >
          <AutoComplete
            placeholder="从列表选择或输入自定义 Client ID"
            options={clientIDList.map((clientID) => ({ value: clientID }))}
            filterOption={(inputValue, option) =>
              option!.value.toUpperCase().indexOf(inputValue.toUpperCase()) !== -1
            }
            allowClear
          />
        </Form.Item>

        <Form.Item
          name="tokensText"
          label="RT 列表"
          rules={[{ required: true, message: '请输入RT' }]}
          extra="每行一个 RT，系统会自动为每个RT生成唯一的32位UUID业务ID，自动去重，跳过已存在的"
        >
          <Input.TextArea 
            rows={12} 
            placeholder="每行一个 RT，系统会自动：&#10;• 为每个RT生成唯一的32位UUID业务ID&#10;• 自动去重（跳过已存在和重复的RT）&#10;• 默认状态：禁用&#10;• 从配置中随机选择代理（如果有配置）&#10;&#10;格式示例：&#10;rt-xxx-1&#10;rt-xxx-2&#10;rt-xxx-3"
            style={{ fontFamily: 'monospace' }}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default BatchImportModal;

