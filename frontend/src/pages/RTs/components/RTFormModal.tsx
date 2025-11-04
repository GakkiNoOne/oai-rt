import React from 'react';
import { Modal, Form, Input, Select, Switch } from 'antd';
import type { FormInstance } from 'antd/es/form';
import type { RT } from '@/api/rts';

interface RTFormModalProps {
  open: boolean;
  editingRT: RT | null;
  form: FormInstance;
  proxyList: string[];
  loadingProxy: boolean;
  onOk: () => void;
  onCancel: () => void;
}

const RTFormModal: React.FC<RTFormModalProps> = ({
  open,
  editingRT,
  form,
  proxyList,
  loadingProxy,
  onOk,
  onCancel,
}) => {
  return (
    <Modal
      title={editingRT ? '编辑RT' : '创建RT'}
      open={open}
      onOk={onOk}
      onCancel={onCancel}
      width={600}
      okText="确定"
      cancelText="取消"
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          enabled: false,
        }}
      >
        <Form.Item
          name="biz_id"
          label="业务ID"
          extra="全表唯一，用于标识此RT（RT每次刷新会变化，业务ID保持不变）"
        >
          <Input placeholder="可选，不输入会自动生成32位UUID（例如：user001、account-01）" />
        </Form.Item>

        <Form.Item
          name="rt_token"
          label="RT"
          rules={[{ required: !editingRT, message: '请输入RT' }]}
        >
          <Input.TextArea 
            rows={3} 
            placeholder={editingRT ? '只读显示' : '请输入RT'}
            disabled={!!editingRT}
            style={editingRT ? { color: '#666', cursor: 'not-allowed' } : {}}
          />
        </Form.Item>

        <Form.Item
          name="proxy"
          label="代理"
        >
          <Select
            placeholder="代理可选，不配置则使用本机ip去请求"
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
          name="tag"
          label="标签"
        >
          <Input placeholder="请输入标签（可选）" allowClear />
        </Form.Item>

        <Form.Item
          name="enabled"
          label="启用状态"
          valuePropName="checked"
        >
          <Switch checkedChildren="启用" unCheckedChildren="禁用" />
        </Form.Item>

        <Form.Item
          name="memo"
          label="备注说明"
        >
          <Input.TextArea 
            rows={3} 
            placeholder="请输入备注说明" 
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default RTFormModal;

