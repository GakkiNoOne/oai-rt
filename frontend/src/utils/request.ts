import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { message } from 'antd';

// 创建axios实例，配置基础URL和默认超时时间
const instance = axios.create({
  baseURL: '/internalweb/v1',
  timeout: 30000, // 默认30秒超时
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器：自动添加token
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：统一错误处理
instance.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error: AxiosError<any>) => {
    // 401 未认证，跳转到登录页
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      message.error('认证已过期，请重新登录');
      window.location.href = '/login';
      return Promise.reject(new Error('未认证'));
    }

    // 超时错误
    if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
      message.error('请求超时，请稍后重试');
      return Promise.reject(new Error('请求超时'));
    }

    // 网络错误
    if (!error.response) {
      message.error('网络错误，请检查网络连接');
      return Promise.reject(new Error('网络错误'));
    }

    // 其他HTTP错误
    const errorData = error.response.data;
    const errorMessage = errorData?.msg || errorData?.error || '请求失败';
    message.error(errorMessage);
    return Promise.reject(new Error(errorMessage));
  }
);

/**
 * GET 请求
 */
function get<T = any>(url: string, options?: { params?: Record<string, any>; timeout?: number }): Promise<T> {
  const config: AxiosRequestConfig = {
    params: options?.params,
    timeout: options?.timeout,
  };
  return instance.get(url, config);
}

/**
 * POST 请求
 */
function post<T = any>(url: string, data?: any, options?: { timeout?: number }): Promise<T> {
  const config: AxiosRequestConfig = {
    timeout: options?.timeout,
  };
  return instance.post(url, data, config);
}

/**
 * PUT 请求
 */
function put<T = any>(url: string, data?: any, options?: { timeout?: number }): Promise<T> {
  const config: AxiosRequestConfig = {
    timeout: options?.timeout,
  };
  return instance.put(url, data, config);
}

/**
 * DELETE 请求
 */
function del<T = any>(url: string, options?: { timeout?: number }): Promise<T> {
  const config: AxiosRequestConfig = {
    timeout: options?.timeout,
  };
  return instance.delete(url, config);
}

// 导出默认对象，包含所有请求方法
const request = {
  get,
  post,
  put,
  delete: del,
};

export default request;
