import React, { useState, useEffect } from 'react';
import { 
  Table, 
  Button, 
  Space, 
  Modal, 
  Form, 
  Tag, 
  Typography, 
  Popconfirm,
  message,
  Tooltip,
  Checkbox,
} from 'antd';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  UploadOutlined,
  SyncOutlined,
  ThunderboltOutlined,
  UserOutlined,
  TeamOutlined,
  CopyOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { rtsApi, type RT, type CreateRTRequest } from '@/api/rts';
import { configsApi } from '@/api/configs';
import RTFormModal from './components/RTFormModal';
import BatchImportModal from './components/BatchImportModal';
import SearchForm from './components/SearchForm';

const { Title } = Typography;

const RTManagement: React.FC = () => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isBatchModalVisible, setIsBatchModalVisible] = useState(false);
  const [editingRT, setEditingRT] = useState<RT | null>(null);
  const [form] = Form.useForm();
  const [batchForm] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [batchLoading, setBatchLoading] = useState(false);
  const [dataSource, setDataSource] = useState<RT[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(100);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [proxyList, setProxyList] = useState<string[]>([]);
  const [loadingProxy, setLoadingProxy] = useState(false);
  const [clientIDList, setClientIDList] = useState<string[]>([]);
  const [loadingClientID, setLoadingClientID] = useState(false);
  
  // 刷新确认弹框状态
  const [isRefreshModalVisible, setIsRefreshModalVisible] = useState(false);
  const [refreshType, setRefreshType] = useState<'single' | 'batch' | 'all'>('single');
  const [refreshTarget, setRefreshTarget] = useState<{ id?: number; biz_id?: string; ids?: number[] } | null>(null);
  const [refreshUserInfo, setRefreshUserInfo] = useState(true);
  const [refreshAccountInfo, setRefreshAccountInfo] = useState(true);
  
  // 搜索条件 - 默认只显示启用的 RT
  const [searchParams, setSearchParams] = useState<{
    biz_id?: string;
    tag?: string;
    email?: string;
    type?: string;
    enabled?: boolean;
    create_date?: string;
  }>({});

  // 加载数据
  const loadData = async () => {
    setLoading(true);
    try {
      const response = await rtsApi.list({
        page,
        page_size: pageSize,
        ...searchParams,
      });
      
      if (response.success && response.data) {
        setDataSource(response.data.items);
        setTotal(response.data.total);
      }
    } catch (error) {
      console.error('加载数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 加载Proxy列表
  const loadProxyList = async () => {
    setLoadingProxy(true);
    try {
      const response = await configsApi.getProxyList();
      if (response.success && response.data) {
        setProxyList(response.data);
      }
    } catch (error) {
      console.error('加载代理列表失败:', error);
    } finally {
      setLoadingProxy(false);
    }
  };

  // 加载 Client ID 列表
  const loadClientIDList = async () => {
    setLoadingClientID(true);
    try {
      const response = await configsApi.getClientIDList();
      if (response.success && response.data) {
        setClientIDList(response.data);
      }
    } catch (error) {
      console.error('加载 Client ID 列表失败:', error);
    } finally {
      setLoadingClientID(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [page, pageSize, searchParams]);

  useEffect(() => {
    loadProxyList();
    loadClientIDList();
  }, []);

  // 处理搜索
  const handleSearch = () => {
    const values = searchForm.getFieldsValue();
    const params: any = {};
    
    if (values.biz_id) params.biz_id = values.biz_id;
    if (values.tag) params.tag = values.tag;
    if (values.email) params.email = values.email;
    if (values.type) params.type = values.type;
    if (values.enabled !== undefined) params.enabled = values.enabled;
    if (values.createDate) {
      params.create_date = values.createDate.format('YYYY-MM-DD');
    }
    
    setSearchParams(params);
    setPage(1); // 重置到第一页
  };

  // 重置搜索
  const handleReset = () => {
    searchForm.resetFields();
    setSearchParams({});
    setPage(1);
  };

  const handleAdd = () => {
    setEditingRT(null);
    form.resetFields();
    form.setFieldsValue({
      enabled: false,
    });
    setIsModalVisible(true);
  };

  const handleEdit = (record: RT) => {
    setEditingRT(record);
    form.setFieldsValue({
      biz_id: record.biz_id,
      rt_token: record.rt,
      proxy: record.proxy,
      client_id: record.client_id,
      tag: record.tag,
      enabled: record.enabled,
      memo: record.memo,
    });
    setIsModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      const response = await rtsApi.delete(id);
      if (response.success) {
        message.success('删除成功');
        loadData();
      }
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();

      if (editingRT) {
        // 编辑 - 只更新允许编辑的字段
        const updateData = {
          biz_id: values.biz_id,
          proxy: values.proxy ?? '',  // 使用 ?? 确保 undefined/null 转为空字符串
          client_id: values.client_id ?? '',
          tag: values.tag ?? '',
          enabled: values.enabled,
          memo: values.memo ?? '',
        };
        console.log('更新数据:', updateData); // 调试日志
        const response = await rtsApi.update(editingRT.id, updateData);
        if (response.success) {
          message.success('更新成功');
          setIsModalVisible(false);
          form.resetFields();
          loadData();
        }
      } else {
        // 新增
        const createData: CreateRTRequest = {
          biz_id: values.biz_id,
          rt_token: values.rt_token,
          proxy: values.proxy,
          client_id: values.client_id,
          tag: values.tag,
          enabled: values.enabled,
          memo: values.memo,
        };
        const response = await rtsApi.create(createData);
        if (response.success) {
          message.success('创建成功');
          setIsModalVisible(false);
          form.resetFields();
          loadData();
        }
      }
    } catch (error) {
      // 错误信息已经在 request.ts 中通过 message.error 显示了
      console.error('操作失败:', error);
    }
  };

  const handleCopyToken = (token: string) => {
    navigator.clipboard.writeText(token);
    message.success('已复制到剪贴板');
  };

  // 渲染截断文本（带Tooltip和可选的点击复制）
  const renderTruncatedText = (text: string | null | undefined, options?: { 
    withCopy?: boolean, 
    codeStyle?: boolean,
    maxLength?: number 
  }) => {
    const { withCopy = false, codeStyle = false, maxLength = 10 } = options || {};
    
    if (!text) {
      return <span style={{ color: '#999' }}>-</span>;
    }

    const truncated = text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
    const needsTruncate = text.length > maxLength;

    const baseStyle = withCopy ? {
      cursor: 'pointer',
      userSelect: 'none' as const,
    } : {};

    const content = codeStyle ? (
      <code style={{ 
        background: '#f5f5f5', 
        padding: '2px 6px', 
        borderRadius: '3px',
        fontSize: '11px',
        ...baseStyle
      }}
      onClick={withCopy ? () => handleCopyToken(text) : undefined}
      >
        {truncated}
      </code>
    ) : (
      <span 
        style={baseStyle}
        onClick={withCopy ? () => handleCopyToken(text) : undefined}
      >
        {truncated}
      </span>
    );

    const textElement = needsTruncate || withCopy ? (
      <Tooltip title={text}>{content}</Tooltip>
    ) : content;

    return textElement;
  };

  // 显示单个刷新确认弹框
  const handleRefresh = (id: number, bizId: string) => {
    setRefreshType('single');
    setRefreshTarget({ id, biz_id: bizId });
    setRefreshUserInfo(true);
    setRefreshAccountInfo(true);
    setIsRefreshModalVisible(true);
  };

  // 执行实际的刷新操作
  const executeRefresh = async () => {
    if (!refreshTarget) return;

    if (refreshType === 'single' && refreshTarget.id) {
      // 单个刷新
      const hide = message.loading(`正在刷新 ${refreshTarget.biz_id}...`, 0);
      try {
        const response = await rtsApi.refresh(refreshTarget.id, refreshUserInfo, refreshAccountInfo);
        hide();
        if (response.success) {
          message.success('刷新成功');
          loadData();
        }
      } catch (error) {
        hide();
        console.error('刷新失败:', error);
      }
    } else if (refreshType === 'batch' && refreshTarget.ids) {
      // 批量刷新
      const hide = message.loading(`正在批量刷新 ${refreshTarget.ids.length} 个RT...`, 0);
      try {
        const response = await rtsApi.batchRefresh(refreshTarget.ids);
        hide();
        if (response.success && response.data) {
          message.success(`批量刷新完成: 成功 ${response.data.success_count} 个, 失败 ${response.data.fail_count} 个`);
          setSelectedRowKeys([]);
          loadData();
        }
      } catch (error) {
        hide();
        console.error('批量刷新失败:', error);
      }
    } else if (refreshType === 'all' && refreshTarget.ids) {
      // 刷新全部
      const hide = message.loading(`正在刷新全部 ${refreshTarget.ids.length} 个RT...`, 0);
      try {
        const response = await rtsApi.batchRefresh(refreshTarget.ids);
        hide();
        if (response.success && response.data) {
          message.success(`刷新全部完成: 成功 ${response.data.success_count} 个, 失败 ${response.data.fail_count} 个`);
          loadData();
        }
      } catch (error) {
        hide();
        console.error('刷新全部失败:', error);
      }
    }

    setIsRefreshModalVisible(false);
    setRefreshTarget(null);
  };

  // 批量导入
  const handleBatchImport = () => {
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  // 批量刷新
  const handleBatchRefresh = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择需要刷新的RT');
      return;
    }

    setRefreshType('batch');
    setRefreshTarget({ ids: selectedRowKeys as number[] });
    setRefreshUserInfo(true);
    setRefreshAccountInfo(true);
    setIsRefreshModalVisible(true);
  };

  // 刷新全部启用的RT
  const handleRefreshAll = () => {
    // 获取所有启用的RT的ID
    const enabledRTs = dataSource.filter(rt => rt.enabled);
    
    if (enabledRTs.length === 0) {
      message.warning('当前页面没有启用的RT可以刷新');
      return;
    }

    const ids = enabledRTs.map(rt => rt.id);
    setRefreshType('all');
    setRefreshTarget({ ids });
    setRefreshUserInfo(true);
    setRefreshAccountInfo(true);
    setIsRefreshModalVisible(true);
  };

  // 刷新用户信息
  const handleRefreshUserInfo = async (id: number, bizId: string) => {
    const hide = message.loading(`正在刷新用户信息 ${bizId}...`, 0);
    try {
      const response = await rtsApi.refreshUserInfo(id);
      hide();
      if (response.success) {
        message.success('刷新用户信息成功');
        loadData();
      }
    } catch (error) {
      hide();
      console.error('刷新用户信息失败:', error);
    }
  };

  // 刷新账号信息
  const handleRefreshAccountInfo = async (id: number, bizId: string) => {
    const hide = message.loading(`正在刷新账号信息 ${bizId}...`, 0);
    try {
      const response = await rtsApi.refreshAccountInfo(id);
      hide();
      if (response.success) {
        message.success('刷新账号信息成功');
        loadData();
      }
    } catch (error) {
      hide();
      console.error('刷新账号信息失败:', error);
    }
  };

  // 批量复制RT
  const handleBatchCopyRT = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择需要复制RT的记录');
      return;
    }

    // 获取选中的RT
    const selectedRTs = dataSource.filter(rt => selectedRowKeys.includes(rt.id));
    
    // 提取RT，过滤掉空值，每行一个RT
    const rtList = selectedRTs
      .filter(rt => rt.rt && rt.rt.trim() !== '')
      .map(rt => rt.rt);
    
    if (rtList.length === 0) {
      message.warning('选中的记录中没有可用的RT');
      return;
    }

    // 每行一个RT（用一个换行符分隔）
    const rtText = rtList.join('\n');
    
    navigator.clipboard.writeText(rtText);
    message.success(`已复制 ${rtList.length} 个RT到剪贴板`);
  };

  // 批量复制AT
  const handleBatchCopyAT = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择需要复制AT的RT');
      return;
    }

    // 获取选中的RT
    const selectedRTs = dataSource.filter(rt => selectedRowKeys.includes(rt.id));
    
    // 提取AT，过滤掉空值，每行一个AT
    const atList = selectedRTs
      .filter(rt => rt.at && rt.at.trim() !== '')
      .map(rt => rt.at);
    
    if (atList.length === 0) {
      message.warning('选中的RT中没有可用的AT');
      return;
    }

    // 每行一个AT（用一个换行符分隔）
    const atText = atList.join('\n');
    
    navigator.clipboard.writeText(atText);
    message.success(`已复制 ${atList.length} 个AT到剪贴板`);
  };

  // 批量删除
  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择需要删除的RT');
      return;
    }

    Modal.confirm({
      title: '确认批量删除',
      content: (
        <div>
          <p style={{ color: '#ff4d4f', marginBottom: 8 }}>
            ⚠️ 警告：此操作不可恢复！
          </p>
          <p>确定要删除选中的 {selectedRowKeys.length} 个RT吗？</p>
        </div>
      ),
      okText: '确定删除',
      cancelText: '取消',
      okButtonProps: { danger: true },
      onOk: async () => {
        const hide = message.loading('正在删除中...', 0);
        try {
          const response = await rtsApi.batchDelete(selectedRowKeys as number[]);
          hide();
          
          if (response.success && response.data) {
            const result = response.data;
            
            if (result.fail_count > 0) {
              message.warning(`删除完成：成功 ${result.success_count} 个，失败 ${result.fail_count} 个`);
            } else {
              message.success(`成功删除 ${result.success_count} 个RT`);
            }
            
            // 清空选中
            setSelectedRowKeys([]);
            // 刷新列表
            loadData();
          }
        } catch (error) {
          hide();
          console.error('批量删除失败:', error);
          message.error('批量删除失败');
        }
      },
    });
  };

  // 处理批量导入
  const handleBatchModalOk = async () => {
    try {
      const values = await batchForm.validateFields();
      const tag = values.tag || '';
      const proxy = values.proxy || '';
      const clientId = values.client_id || '';
      const tokensText = values.tokensText;
      
      // 按行分割，过滤空行
      const rtTokens = tokensText
        .split('\n')
        .map((line: string) => line.trim())
        .filter((line: string) => line !== '');
      
      if (rtTokens.length === 0) {
        message.error('请输入至少一个RT');
        return;
      }
      
      setBatchLoading(true);
      // 注意：batchName参数已弃用，每个RT会生成唯一的随机8位ID
      const response = await rtsApi.batchCreate('', tag, rtTokens, clientId, proxy);
      
      if (response.success && response.data) {
        const result = response.data;
        message.success(response.msg);
        
        // 显示详细结果
        Modal.info({
          title: '批量导入结果',
          content: (
            <div>
              <p>总数量：{result.total_count}</p>
              <p style={{ color: '#52c41a' }}>✅ 成功导入：{result.success_count}</p>
              {result.fail_count > 0 && (
                <p style={{ color: '#faad14' }}>
                  ⚠️ 跳过：{result.fail_count}（重复或已存在）
                </p>
              )}
            </div>
          ),
        });
        
        setIsBatchModalVisible(false);
        batchForm.resetFields();
        loadData();
      }
    } catch (error) {
      console.error('批量导入失败:', error);
    } finally {
      setBatchLoading(false);
    }
  };

  // 切换状态
  const handleToggleStatus = (record: RT) => {
    const newStatus = !record.enabled;
    const statusText = newStatus ? '启用' : '禁用';
    
    Modal.confirm({
      title: `确认${statusText}`,
      content: `确定要${statusText}RT "${record.biz_id}" 吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          const response = await rtsApi.update(record.id, {
            enabled: newStatus,
          });
          if (response.success) {
            message.success(`${statusText}成功`);
            loadData();
          }
        } catch (error) {
          console.error(`${statusText}失败:`, error);
        }
      },
    });
  };

  const getStatusTag = (enabled: boolean, record: RT) => {
    return (
      <Tag 
        color={enabled ? "green" : "red"}
        style={{ cursor: 'pointer' }}
        onClick={() => handleToggleStatus(record)}
      >
        {enabled ? '启用' : '禁用'}
      </Tag>
    );
  };

  const columns: ColumnsType<RT> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 70,
      fixed: 'left' as const,
    },
    {
      title: '创建时间',
      dataIndex: 'create_time',
      key: 'create_time',
      width: 140,
      render: (date) => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: '业务ID',
      dataIndex: 'biz_id',
      key: 'biz_id',
      width: 140,
      render: (text) => renderTruncatedText(text, { withCopy: true, maxLength: 10 }),
    },
    {
      title: '名称',
      dataIndex: 'user_name',
      key: 'user_name',
      width: 120,
      render: (text) => text || <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 120,
      render: (text) => text || <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (text) => text || <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: '上次刷新',
      dataIndex: 'last_refresh_time',
      key: 'last_refresh_time',
      width: 140,
      render: (date) => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: 'RT',
      dataIndex: 'rt',
      key: 'rt',
      width: 140,
      render: (token) => renderTruncatedText(token, { withCopy: true, codeStyle: true, maxLength: 10 }),
    },
    {
      title: 'AT',
      dataIndex: 'at',
      key: 'at',
      width: 140,
      render: (token) => {
        if (!token) {
          return <span style={{ color: '#999' }}>-</span>;
        }
        const truncated = token.length > 10 ? token.substring(0, 10) + '...' : token;
        return (
          <code 
            style={{ 
              background: '#f5f5f5', 
              padding: '2px 6px', 
              borderRadius: '3px',
              fontSize: '11px',
              cursor: 'pointer',
              userSelect: 'none',
            }}
            onClick={() => handleCopyToken(token)}
          >
            {truncated}
          </code>
        );
      },
    },
    {
      title: 'Tag',
      dataIndex: 'tag',
      key: 'tag',
      width: 100,
      render: (tag) => {
        if (!tag) {
          return <span style={{ color: '#999' }}>-</span>;
        }
        return <Tag color="blue">{tag}</Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 80,
      render: (enabled, record) => getStatusTag(enabled, record),
    },
    {
      title: 'Proxy',
      dataIndex: 'proxy',
      key: 'proxy',
      width: 160,
      ellipsis: true,
      render: (proxy) => proxy || <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: 'ClientId',
      dataIndex: 'client_id',
      key: 'client_id',
      width: 200,
      ellipsis: true,
      render: (clientId) => clientId || <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: '备注',
      dataIndex: 'memo',
      key: 'memo',
      width: 150,
      ellipsis: true,
      render: (memo) => memo || <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: '上一个RT',
      dataIndex: 'last_rt',
      key: 'last_rt',
      width: 140,
      render: (lastRt) => renderTruncatedText(lastRt, { withCopy: true, codeStyle: true, maxLength: 10 }),
    },
    {
      title: '刷新结果',
      dataIndex: 'refresh_result',
      key: 'refresh_result',
      width: 120,
      render: (result) => {
        if (!result) {
          return <span style={{ color: '#999' }}>-</span>;
        }
        const truncated = result.length > 10 ? result.substring(0, 10) + '...' : result;
        return (
          <span 
            style={{ 
              cursor: 'pointer',
              userSelect: 'none',
            }}
            onClick={() => handleCopyToken(result)}
          >
            {truncated}
          </span>
        );
      },
    },
    {
      title: '用户信息',
      dataIndex: 'user_info',
      key: 'user_info',
      width: 120,
      render: (info) => {
        if (!info) {
          return <span style={{ color: '#999' }}>-</span>;
        }
        const truncated = info.length > 10 ? info.substring(0, 10) + '...' : info;
        return (
          <span 
            style={{ 
              cursor: 'pointer',
              userSelect: 'none',
            }}
            onClick={() => handleCopyToken(info)}
          >
            {truncated}
          </span>
        );
      },
    },
    {
      title: '账号信息',
      dataIndex: 'account_info',
      key: 'account_info',
      width: 120,
      render: (info) => {
        if (!info) {
          return <span style={{ color: '#999' }}>-</span>;
        }
        const truncated = info.length > 10 ? info.substring(0, 10) + '...' : info;
        return (
          <span 
            style={{ 
              cursor: 'pointer',
              userSelect: 'none',
            }}
            onClick={() => handleCopyToken(info)}
          >
            {truncated}
          </span>
        );
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      fixed: 'right' as const,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="刷新RT">
            <Button 
              type="link" 
              size="small"
              icon={<SyncOutlined />}
              onClick={() => handleRefresh(record.id, record.biz_id)}
            />
          </Tooltip>
          <Tooltip title="刷新用户信息">
            <Button 
              type="link" 
              size="small"
              icon={<UserOutlined />}
              onClick={() => handleRefreshUserInfo(record.id, record.biz_id)}
            />
          </Tooltip>
          <Tooltip title="刷新账号信息">
            <Button 
              type="link" 
              size="small"
              icon={<TeamOutlined />}
              onClick={() => handleRefreshAccountInfo(record.id, record.biz_id)}
            />
          </Tooltip>
          <Tooltip title="编辑">
            <Button 
              type="link" 
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          <Popconfirm
            title="确定要删除这个RT吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button 
                type="link" 
                size="small"
                danger 
                icon={<DeleteOutlined />}
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      {/* 标题 */}
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: 16
      }}>
        <div>
          <Title level={2} style={{ marginBottom: 0 }}>RT管理</Title>
        </div>
        <Space>
          <Tooltip title="刷新数据库中所有启用的RT">
            <Button 
              icon={<ThunderboltOutlined />}
              onClick={handleRefreshAll}
              type="dashed"
            >
              刷新全部
            </Button>
          </Tooltip>
          <Button 
            icon={<UploadOutlined />}
            onClick={handleBatchImport}
          >
            批量导入
          </Button>
          <Button 
            type="primary" 
            icon={<PlusOutlined />}
            onClick={handleAdd}
          >
            创建RT
          </Button>
        </Space>
      </div>

      {/* 搜索表单 */}
      <SearchForm
        form={searchForm}
        onSearch={handleSearch}
        onReset={handleReset}
      />

      {/* 批量操作区域 - 默认显示 */}
      <div style={{ 
        marginBottom: 16, 
        padding: '12px 16px',
        background: selectedRowKeys.length > 0 ? '#e6f7ff' : '#f5f5f5',
        border: selectedRowKeys.length > 0 ? '1px solid #91d5ff' : '1px solid #d9d9d9',
        borderRadius: '4px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <span style={{ color: selectedRowKeys.length > 0 ? '#1890ff' : '#999' }}>
          已选择 <strong>{selectedRowKeys.length}</strong> 项
        </span>
        <Space>
          <Button
            size="small"
            icon={<SyncOutlined />}
            onClick={handleBatchRefresh}
            disabled={selectedRowKeys.length === 0}
          >
            批量刷新
          </Button>
          <Button
            size="small"
            icon={<CopyOutlined />}
            onClick={handleBatchCopyRT}
            disabled={selectedRowKeys.length === 0}
          >
            复制RT
          </Button>
          <Button
            size="small"
            icon={<CopyOutlined />}
            onClick={handleBatchCopyAT}
            disabled={selectedRowKeys.length === 0}
          >
            复制AT
          </Button>
          <Button
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={handleBatchDelete}
            disabled={selectedRowKeys.length === 0}
          >
            批量删除
          </Button>
          {selectedRowKeys.length > 0 && (
            <Button 
              size="small" 
              type="link"
              onClick={() => setSelectedRowKeys([])}
            >
              取消选择
            </Button>
          )}
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={dataSource}
        rowKey="id"
        loading={loading}
        scroll={{ x: 2340 }}
        rowSelection={{
          selectedRowKeys,
          onChange: (selectedRowKeys) => {
            setSelectedRowKeys(selectedRowKeys);
          },
          preserveSelectedRowKeys: true,
        }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) => 
            `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
      />

      {/* 创建/编辑弹窗 */}
      <RTFormModal
        open={isModalVisible}
        editingRT={editingRT}
        form={form}
        proxyList={proxyList}
        loadingProxy={loadingProxy}
        clientIDList={clientIDList}
        loadingClientID={loadingClientID}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
      />

      {/* 批量导入弹窗 */}
      <BatchImportModal
        open={isBatchModalVisible}
        form={batchForm}
        loading={batchLoading}
        proxyList={proxyList}
        loadingProxy={loadingProxy}
        clientIDList={clientIDList}
        loadingClientID={loadingClientID}
        onOk={handleBatchModalOk}
        onCancel={() => {
          setIsBatchModalVisible(false);
          batchForm.resetFields();
        }}
      />

      {/* 刷新确认弹窗 */}
      <Modal
        title={
          refreshType === 'single' ? '刷新RT' :
          refreshType === 'batch' ? '批量刷新RT' :
          '刷新全部RT'
        }
        open={isRefreshModalVisible}
        onOk={executeRefresh}
        onCancel={() => {
          setIsRefreshModalVisible(false);
          setRefreshTarget(null);
        }}
        okText="确定刷新"
        cancelText="取消"
        width={450}
      >
        <div style={{ marginBottom: 16 }}>
          {refreshType === 'single' && refreshTarget?.biz_id && (
            <p>确定要刷新 <strong>{refreshTarget.biz_id}</strong> 吗？</p>
          )}
          {refreshType === 'batch' && refreshTarget?.ids && (
            <p>确定要刷新选中的 <strong>{refreshTarget.ids.length}</strong> 个RT吗？</p>
          )}
          {refreshType === 'all' && refreshTarget?.ids && (
            <p>确定要刷新当前页面所有启用的 <strong>{refreshTarget.ids.length}</strong> 个RT吗？</p>
          )}
        </div>
        
        {/* 只在单个刷新时显示刷新选项 */}
        {refreshType === 'single' && (
          <div style={{ padding: '16px', background: '#f5f5f5', borderRadius: '4px' }}>
            <p style={{ marginBottom: 12, fontWeight: 500 }}>刷新选项：</p>
            <div style={{ marginBottom: 8 }}>
              <Checkbox
                checked={refreshUserInfo}
                onChange={(e) => setRefreshUserInfo(e.target.checked)}
              >
                同时刷新用户信息（邮箱、名称）
              </Checkbox>
            </div>
            <div>
              <Checkbox
                checked={refreshAccountInfo}
                onChange={(e) => setRefreshAccountInfo(e.target.checked)}
              >
                同时刷新账号信息（类型）
              </Checkbox>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default RTManagement;
